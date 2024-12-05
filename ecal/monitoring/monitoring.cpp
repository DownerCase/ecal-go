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
  const auto clients     = toCTypes<CClientMon>(monitoring.clients);
  const auto servers     = toCTypes<CServerMon>(monitoring.server);
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
