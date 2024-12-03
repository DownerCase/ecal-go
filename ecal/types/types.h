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

struct CTopicMon {
  const char *unit_name;
  const char *topic_id;
  const char *topic_name;
  const char *direction;
  struct CDatatype datatype;
  int64_t data_clock;
  int64_t data_freq;
  int32_t registration_clock;
  int32_t topic_size;
  int32_t connections_local;
  int32_t connections_external;
  int32_t message_drops;
};

struct CProcessMon {
  const char *host_name;
  const char *host_group;
  int32_t registration_clock;
  int32_t pid;
  const char *process_name;
  const char *unit_name;
  const char *process_parameters;
  int32_t state_severity;
  int32_t state_severity_level;
  const char *state_info;
  const char *components;
  const char *runtime;
};

#ifdef __cplusplus
}
#endif

#endif
