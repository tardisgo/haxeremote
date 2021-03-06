#include <stdio.h>

extern "C" const char *hxRunLibrary();
extern "C" void hxcpp_set_top_of_stack();

extern "C" int hxmain()	
{
  hxcpp_set_top_of_stack();
  const char *err = hxRunLibrary();
	if (err) {
		// Unhandled exceptions ...
		fprintf(stderr,"Error %s\n", err );
		return -1;
	}
	return 0;
}

