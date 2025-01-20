#ifndef ECAL_GO_TYPES
#define ECAL_GO_TYPES

#include <stddef.h>
#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

struct CDatatype {
  const char *name;
  const char *encoding;
  const void *descriptor;
  int descriptor_len;
  char _PADDING[4];
};

#ifdef __cplusplus
}
#endif

#endif
