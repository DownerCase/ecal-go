#ifndef ECAL_GO_TYPES
#define ECAL_GO_TYPES

#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

struct CEntityId {
  const char *entity_id;
  int entity_id_len;
  int32_t process_id;
  const char *host_name;
  int host_name_len;
};

struct CTopicId {
  struct CEntityId topic_id;
  const char *topic_name;
  int topic_name_len;
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
