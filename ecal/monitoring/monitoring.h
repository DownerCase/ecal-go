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

struct CMonitoring {
  const struct CTopicMon *publishers;
  size_t publishers_len;
  const struct CTopicMon *subscribers;
  size_t subscribers_len;
  const struct CProcessMon *processes;
  size_t processes_len;
};

void GetMonitoring(uintptr_t handle, unsigned int entities);

#ifdef __cplusplus
}
#endif

#endif
