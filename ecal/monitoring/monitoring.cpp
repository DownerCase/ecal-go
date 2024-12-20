#include "monitoring.h"

#include "types.h"
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

void GetMonitoring(uintptr_t handle, unsigned int entities) {

  eCAL::Monitoring::SMonitoring monitoring{};
  eCAL::Monitoring::GetMonitoring(monitoring, entities);
  const auto publishers  = toCTypes<CTopicMon>(monitoring.publisher);
  const auto subscribers = toCTypes<CTopicMon>(monitoring.subscriber);
  const auto processes   = toCTypes<CProcessMon>(monitoring.processes);

  std::vector<CClientMon> clients{};
  std::vector<CServerMon> servers{};
  std::vector<std::vector<CMethodMon>> serviceMethods{};

  for (const auto &client : monitoring.clients) {
    clients.emplace_back(toCType(client));
    auto &cclient = clients.back();
    serviceMethods.emplace_back(toCTypes<CMethodMon>(client.methods));
    cclient.base.methods_len = serviceMethods.back().size();
    cclient.base.methods     = serviceMethods.back().data();
  }
  for (const auto &server : monitoring.server) {
    servers.emplace_back(toCType(server));
    auto &cserver = servers.back();
    serviceMethods.emplace_back(toCTypes<CMethodMon>(server.methods));
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
