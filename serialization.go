// package haxeremote provides Serialization and Unserializtion of data in Haxe format.
/* Haxe serialization prefixes :
a : array 					- DONE
b : hash
c : class
d : Float 					- DONE
e : reserved (float exp) 	- DONE
f : false 					- DONE
g : object end
h : array/list/hash end 	- DONE for array
i : Int 					- DONE
j : enum (by index)
k : NaN 					- DONE
l : list
m : -Inf 					- DONE
n : null 					- DONE
o : object
p : +Inf 					- DONE
q : haxe.ds.IntMap
r : reference
s : bytes (base64) 			- DONE
t : true					- DONE
u : array nulls				- DONE
v : date
w : enum
x : exception
y : urlencoded string		- DONE
z : zero					- DONE
A : Class<Dynamic>
B : Enum<Dynamic>
M : haxe.ds.ObjectMap
C : custom
*/
package haxeremote

import (
	"encoding/base64"
	"errors"
	"fmt"
	"math"
	"net/url"
	"strconv"
)

// Unserialize the in the Haxe format from a given buffer.
func Unserialize(buf []byte) (data interface{}, remains []byte, err error) {
	//fmt.Println("DEBUG Unserialize:", string(buf))
	if len(buf) == 0 {
		return nil, nil, nil
	}
	code := buf[0]
	remains = buf[1:]
	switch code {
	case 'a': // Array
		dataArray := []interface{}{}
		var item interface{}
	arrayLoop:
		item, remains, err = Unserialize(remains)
		if err != nil {
			return nil, nil, err
		}
		dataArray = append(dataArray, item)
	arrayLoopEnd:
		if len(remains) > 0 && err == nil {
			switch remains[0] {
			default:
				goto arrayLoop
			case 'u': // some number of null/nil entries
				remains = remains[1:]
				for i := range remains {
					switch remains[i] {
					case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
					default:
						strInt := remains[:i]
						remains = remains[i:]
						numNil := 0
						numNil, err = strconv.Atoi(string(strInt))
						for j := 0; j < numNil; j++ {
							dataArray = append(dataArray, nil)
						}
						goto arrayLoopEnd
					}
				}
				return nil, nil,
					errors.New("invalid u item in Haxe array serialization: " + string(remains))
			case 'h': // end of array
				remains = remains[1:]
			}
		}
		data = dataArray

	case 'i': // Integer
		for i := range remains {
			switch remains[i] {
			case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			default:
				strInt := remains[:i]
				remains = remains[i:]
				data, err = strconv.Atoi(string(strInt))
				goto done
			}
		}
		data, err = strconv.Atoi(string(remains))
		remains = []byte{}

	case 'y': // String
		for i, j := range remains {
			switch j {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				// NoOp
			case ':':
				strInt := string(remains[:i])
				remains = remains[i+1:]
				var length int
				length, err = strconv.Atoi(strInt)
				//fmt.Printf("DEBUG string len: %s, %d, %v\n", strInt, length, err)
				if err == nil {
					raw := string(remains[:length])
					clean, err := url.QueryUnescape(raw)
					if err == nil {
						data = clean
					} else {
						data = raw
					}
					remains = remains[length:]
					//fmt.Printf("DEBUG string decoded, remaining: %s, %s\n", data, remains)
				}
				goto done
			default:
				err = errors.New("unrecognised string length: " + string(remains))
			}
		}

	case 'd': // Float
		for i := range remains {
			switch remains[i] {
			case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '.', 'e':
			default:
				strFloat := string(remains[:i])
				remains = remains[i:]
				data, err = strconv.ParseFloat(strFloat, 64)
				goto done
			}
		}
		data, err = strconv.ParseFloat(string(remains), 64)
		remains = []byte{}

		// the single letter values
	case 'n':
		data = nil
	case 't':
		data = true
	case 'f':
		data = false
	case 'k':
		data = math.NaN()
	case 'p':
		data = math.Inf(+1)
	case 'm':
		data = math.Inf(-1)
	case 'z': // TODO should zero be floating point?
		data = 0

	case 's': // haxe.io.Bytes
		for i, j := range remains {
			switch j {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				// NoOp
			case ':':
				strInt := string(remains[:i])
				remains = remains[i+1:]
				var length int
				length, err = strconv.Atoi(strInt)
				if err == nil {
					raw := string(remains[:length])
					//println("DEBUG raw=" + raw)
					data, err = base64.StdEncoding.DecodeString(raw)
					//fmt.Printf("DEBUG data %v:%T err=%v\n", data, data, err)
					remains = remains[length:]
				}
				goto done
			default:
				err = errors.New("unrecognised string length: " + string(remains))
			}
		}

		// TODO many more letters !

	default:
		err = errors.New("unhandled Haxe serialization code: " + string(code))
	}
done:
	return
}

// Serialize the given data into a string in Haxe serialization format.
func Serialize(data interface{}) string {
	if data == nil {
		return "n"
	}
	switch data.(type) {
	case int:
		if data.(int) == 0 {
			return "z"
		}
		return fmt.Sprintf("i%d", data.(int))

	case string:
		result := url.QueryEscape(data.(string))
		return fmt.Sprintf("y%d:%s", len(result), result)

	case bool:
		if data.(bool) {
			return "t"
		}
		return "f"

	case float64: // NOTE no special processing of 0=>z for Float
		val := data.(float64)
		if math.IsInf(val, -1) {
			return "m"
		}
		if math.IsInf(val, +1) {
			return "p"
		}
		if math.IsNaN(val) {
			return "k"
		}
		return "d" + strconv.FormatFloat(val, 'g', -1, 64)

	case []interface{}:
		ret := "a"
		nilCount := 0
		for _, a := range data.([]interface{}) {
			if a == nil {
				nilCount++
			} else {
				if nilCount > 0 {
					ret += fmt.Sprintf("u%d", nilCount)
					nilCount = 0
				}
				ret += Serialize(a)
			}
		}
		if nilCount > 0 {
			ret += fmt.Sprintf("u%d", nilCount)
		}
		return ret + "h"

	case []byte:
		strForm := base64.StdEncoding.EncodeToString(data.([]byte))
		return fmt.Sprintf("s%d:%s", len(strForm), strForm)

	default:
		panic(fmt.Sprintf("unhandled type to be serialized for Haxe %v : %T", data, data))
	}
}
