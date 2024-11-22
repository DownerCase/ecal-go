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
CTopicMon toCTopicMon(const eCAL::Monitoring::STopicMon &topic) {
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

std::vector<CTopicMon>
toCTopics(const std::vector<eCAL::Monitoring::STopicMon> &topics) {
  std::vector<CTopicMon> ctopics{};

  for (const auto &topic : topics) {
    ctopics.emplace_back(toCTopicMon(topic));
  }

  return ctopics;
}
} // namespace

void GetMonitoring(uintptr_t handle, unsigned int entities) {

  eCAL::Monitoring::SMonitoring monitoring{};
  eCAL::Monitoring::GetMonitoring(monitoring, entities);
  const auto publishers  = toCTopics(monitoring.publisher);
  const auto subscribers = toCTopics(monitoring.subscriber);
  CMonitoring cmon{};
  cmon.publishers      = publishers.data();
  cmon.publishers_len  = publishers.size();
  cmon.subscribers     = subscribers.data();
  cmon.subscribers_len = subscribers.size();
  goCopyMonitoring(handle, &cmon);
}
