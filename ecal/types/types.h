#ifndef ECAL_GO_TYPES
#define ECAL_GO_TYPES

#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

struct CEntityId {
  const char *entity_id;
  const char *host_name;
  int32_t process_id;
};

struct CTopicId {
  struct CEntityId topic_id;
  const char *topic_name;
};

struct CDatatype {
  const char *name;
  const char *encoding;
  const void *descriptor;
  int descriptor_len;
};

#ifdef __cplusplus
}
#endif

#endif
