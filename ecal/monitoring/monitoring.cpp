#include "monitoring.h"

#include "types.hpp"

#include <ecal/ecal_monitoring.h>
#include <ecal/types/monitoring.h>

extern "C" {
extern void goCopyMonitoring(uintptr_t handle, CMonitoring *);
}

// Ensure exposed constants are correct
static_assert(monitoring_none == eCAL::Monitoring::Entity::None);
static_assert(monitoring_publisher == eCAL::Monitoring::Entity::Publisher);
static_assert(monitoring_subscriber == eCAL::Monitoring::Entity::Subscriber);
static_assert(monitoring_server == eCAL::Monitoring::Entity::Server);
static_assert(monitoring_client == eCAL::Monitoring::Entity::Client);
static_assert(monitoring_process == eCAL::Monitoring::Entity::Process);
static_assert(monitoring_host == eCAL::Monitoring::Entity::Host);
static_assert(monitoring_all == eCAL::Monitoring::Entity::All);

namespace {

// Makes uses of copy elision
CTopicMon toCType(const eCAL::Monitoring::STopicMon &topic) {
  return {
      topic.rclock,
      topic.uname.c_str(),
      topic.tid.c_str(),
      topic.tname.c_str(),
      topic.direction.c_str(),
      toCDataType(topic.tdatatype),
      topic.tsize,
      topic.connections_loc,
      topic.connections_ext,
      topic.message_drops,
      topic.dclock,
      topic.dfreq
  };
}

CProcessMon toCType(const eCAL::Monitoring::SProcessMon &proc) {
  return {
      proc.hname.c_str(),
      proc.hgname.c_str(),
      proc.rclock,
      proc.pid,
      proc.pname.c_str(),
      proc.uname.c_str(),
      proc.pparam.c_str(),
      proc.state_severity,
      proc.state_severity_level,
      proc.state_info.c_str(),
      proc.component_init_info.c_str(),
      proc.ecal_runtime_version.c_str()
  };
}

template <typename CType, typename EcalType>
std::vector<CType> toCTypes(const std::vector<EcalType> &ecaltypes) {
  std::vector<CType> ctypes{};
  ctypes.reserve(ecaltypes.size());
  for (const auto &ecaltype : ecaltypes) {
    ctypes.emplace_back(toCType(ecaltype));
  }
  return ctypes;
}

} // namespace

void GetMonitoring(uintptr_t handle, unsigned int entities) {

  eCAL::Monitoring::SMonitoring monitoring{};
  eCAL::Monitoring::GetMonitoring(monitoring, entities);
  const auto publishers  = toCTypes<CTopicMon>(monitoring.publisher);
  const auto subscribers = toCTypes<CTopicMon>(monitoring.subscriber);
  const auto processes   = toCTypes<CProcessMon>(monitoring.processes);
  CMonitoring cmon{};
  cmon.publishers      = publishers.data();
  cmon.publishers_len  = publishers.size();
  cmon.subscribers     = subscribers.data();
  cmon.subscribers_len = subscribers.size();
  cmon.processes       = processes.data();
  cmon.processes_len   = processes.size();
  goCopyMonitoring(handle, &cmon);
}
