#ifndef ECAL_GO_TYPES
#define ECAL_GO_TYPES

#include <stddef.h>
#include <stdint.h>

#include <ecal/ecal_log_level.h>

#ifdef __cplusplus
extern "C" {
#endif

struct CEntityId {
  uint64_t entity_id;
  const char *host_name;
  int32_t process_id;
  char PADDING[4];
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
  char _PADDING[4];
};

struct CTopicMon {
  const char *unit_name;
  const char *host_name;
  uint64_t topic_id;
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
  char _PADDING[4];
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

struct CMethodMon {
  const char *name;
  struct CDatatype req_datatype;
  struct CDatatype resp_datatype;
  int64_t call_count;
};

struct CServiceCommon {
  const char *name;
  uint64_t id;
  const struct CMethodMon *methods;
  size_t methods_len;
  const char *host_name;
  const char *process_name;
  const char *unit_name;
  int32_t registration_clock;
  int32_t pid;
  uint32_t protocol_version;
  char _PADDING[4];
};

struct CServerMon {
  struct CServiceCommon base;
  uint32_t port_v0;
  uint32_t port_v1;
};

struct CClientMon {
  struct CServiceCommon base;
};

struct CLogMessage {
  int64_t time;
  const char *host_name;
  const char *process_name;
  const char *unit_name;
  const char *content;
  int32_t pid;
  enum eCAL_Logging_eLogLevel level;
};

#ifdef __cplusplus
}
#endif

#endif
