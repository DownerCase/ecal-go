#ifndef ECAL_GO_MONITORING
#define ECAL_GO_MONITORING

#include <stddef.h>
#include <stdint.h>

#include "types.h"

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
  int32_t registration_clock;
  const char *unit_name;
  const char *topic_id;
  const char *topic_name;
  const char *direction;
  struct CDatatype datatype;
  int32_t topic_size;
  int32_t connections_local;
  int32_t connections_external;
  int32_t message_drops;
  int64_t data_clock;
  int64_t data_freq;
};

struct CMonitoring {
  const struct CTopicMon *publishers;
  size_t publishers_len;
  const struct CTopicMon *subscribers;
  size_t subscribers_len;
};

void GetMonitoring(uintptr_t handle, unsigned int entities);

#ifdef __cplusplus
}
#endif

#endif
