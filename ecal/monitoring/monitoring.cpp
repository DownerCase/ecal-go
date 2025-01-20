#include "monitoring.h"

#include "types.h"
#include "types.hpp"

#include <ecal/ecal_monitoring.h>
#include <ecal/types/monitoring.h>

extern "C" {
extern void goCopyMonitoring(uintptr_t handle, CMonitoring *);
}

namespace {
CTopicMon toCTopicMon(const eCAL::Monitoring::STopicMon &topic) {
  return {
      topic.uname.c_str(),
      topic.hname.c_str(),
      topic.tid,
      topic.tname.c_str(),
      topic.direction.c_str(),
      toCType(topic.tdatatype),
      topic.dclock,
      topic.dfreq,
      topic.rclock,
      topic.tsize,
      topic.connections_loc,
      topic.connections_ext,
      topic.message_drops,
      {}
  };
}

CProcessMon toCProcessMon(const eCAL::Monitoring::SProcessMon &proc) {
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

template <class T> CServiceCommon toCServiceCommon(const T &t) {
  return {
      t.sname.c_str(),
      t.sid,
      nullptr, // Methods are filled in a separate pass
      t.methods.size(),
      t.hname.c_str(),
      t.pname.c_str(),
      t.uname.c_str(),
      t.rclock,
      t.pid,
      t.version,
      {}
  };
}
CClientMon toCClientMon(const eCAL::Monitoring::SClientMon &client) {
  return {toCServiceCommon(client)};
}

CServerMon toCServerMon(const eCAL::Monitoring::SServerMon &server) {
  return {toCServiceCommon(server), server.tcp_port_v0, server.tcp_port_v1};
}

CMethodMon toCMethodMon(const eCAL::Monitoring::SMethodMon &method) {
  return {
      method.mname.c_str(),
      toCType(method.req_datatype),
      toCType(method.resp_datatype),
      method.call_count
  };
}

} // namespace

// Ensure exposed constants are correct
static_assert(monitoring_none == eCAL::Monitoring::Entity::None);
static_assert(monitoring_publisher == eCAL::Monitoring::Entity::Publisher);
static_assert(monitoring_subscriber == eCAL::Monitoring::Entity::Subscriber);
static_assert(monitoring_server == eCAL::Monitoring::Entity::Server);
static_assert(monitoring_client == eCAL::Monitoring::Entity::Client);
static_assert(monitoring_process == eCAL::Monitoring::Entity::Process);
static_assert(monitoring_host == eCAL::Monitoring::Entity::Host);
static_assert(monitoring_all == eCAL::Monitoring::Entity::All);

void GetMonitoring(uintptr_t handle, unsigned int entities) {

  eCAL::Monitoring::SMonitoring monitoring{};
  eCAL::Monitoring::GetMonitoring(monitoring, entities);
  const auto publishers =
      containerTo<std::vector>(monitoring.publisher, toCTopicMon);
  const auto subscribers =
      containerTo<std::vector>(monitoring.subscriber, toCTopicMon);
  const auto processes =
      containerTo<std::vector>(monitoring.processes, toCProcessMon);

  std::vector<CClientMon> clients{};
  std::vector<CServerMon> servers{};
  std::vector<std::vector<CMethodMon>> serviceMethods{};

  for (const auto &client : monitoring.clients) {
    clients.emplace_back(toCClientMon(client));
    auto &cclient = clients.back();
    serviceMethods.emplace_back(
        containerTo<std::vector>(client.methods, toCMethodMon)
    );
    cclient.base.methods_len = serviceMethods.back().size();
    cclient.base.methods     = serviceMethods.back().data();
  }
  for (const auto &server : monitoring.server) {
    servers.emplace_back(toCServerMon(server));
    auto &cserver = servers.back();
    serviceMethods.emplace_back(
        containerTo<std::vector>(server.methods, toCMethodMon)
    );
    cserver.base.methods_len = serviceMethods.back().size();
    cserver.base.methods     = serviceMethods.back().data();
  }

  CMonitoring cmon{};
  cmon.publishers      = publishers.data();
  cmon.publishers_len  = publishers.size();
  cmon.subscribers     = subscribers.data();
  cmon.subscribers_len = subscribers.size();
  cmon.processes       = processes.data();
  cmon.processes_len   = processes.size();
  cmon.clients         = clients.data();
  cmon.clients_len     = clients.size();
  cmon.servers         = servers.data();
  cmon.servers_len     = servers.size();
  goCopyMonitoring(handle, &cmon);
}
