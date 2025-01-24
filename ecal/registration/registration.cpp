#include "registration.h"

#include <climits>
#include <ecal/registration.h>

extern "C" {
extern void goTopicEventCallback(uintptr_t handle, CTopicId id, uint8_t event);
}

namespace {
CTopicId toCTopicId(const eCAL::Registration::STopicId &id) {
  return {
      {
          id.topic_id.entity_id,
          id.topic_id.host_name.data(),
          id.topic_id.process_id,
          {},
      },
      id.topic_name.data()
  };
}
} // namespace

size_t AddPublisherEventCallback(uintptr_t handle) {
  const auto callback_adapter = [handle](
                                    const eCAL::Registration::STopicId &id,
                                    eCAL::Registration::RegistrationEventType
                                        event
                                ) {
    goTopicEventCallback(handle, toCTopicId(id), static_cast<uint8_t>(event));
  };
  return eCAL::Registration::AddPublisherEventCallback(callback_adapter);
}

void RemPublisherEventCallback(size_t handle) {
  eCAL::Registration::RemPublisherEventCallback(handle);
}

size_t AddSubscriberEventCallback(uintptr_t handle) {
  const auto callback_adapter = [handle](
                                    const eCAL::Registration::STopicId &id,
                                    eCAL::Registration::RegistrationEventType
                                        event
                                ) {
    goTopicEventCallback(handle, toCTopicId(id), static_cast<uint8_t>(event));
  };
  return eCAL::Registration::AddSubscriberEventCallback(callback_adapter);
}

void RemSubscriberEventCallback(size_t handle) {
  eCAL::Registration::RemSubscriberEventCallback(handle);
}
