#ifndef ECAL_GO_MONITORING
#define ECAL_GO_MONITORING

#include <stddef.h>
#include <stdint.h>

#include "cgo_types.h"

#ifdef __cplusplus
extern "C" {
#endif

enum eCAL_Monitoring_eEntity {
  monitoring_none       = 0x000,
  monitoring_publisher  = 0x001,
  monitoring_subscriber = 0x002,
  monitoring_server     = 0x004,
  monitoring_client     = 0x008,
  monitoring_process    = 0x010,
  monitoring_host       = 0x020,
  monitoring_all        = monitoring_publisher | monitoring_subscriber |
                   monitoring_server | monitoring_client | monitoring_process |
                   monitoring_host,
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
  const char *shm_domain;
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

struct CMonitoring {
  const struct CTopicMon *publishers;
  size_t publishers_len;
  const struct CTopicMon *subscribers;
  size_t subscribers_len;
  const struct CProcessMon *processes;
  size_t processes_len;
  const struct CClientMon *clients;
  size_t clients_len;
  const struct CServerMon *servers;
  size_t servers_len;
};

void GetMonitoring(uintptr_t handle, unsigned int entities);

#ifdef __cplusplus
}
#endif

#endif
