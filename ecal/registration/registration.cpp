#include "registration.h"

#include "types.hpp"

#include <climits>
#include <ecal/ecal_registration.h>

extern "C" {
extern void goTopicEventCallback(uintptr_t handle, CTopicId id, uint8_t event);
}

namespace {
int safe_len(size_t str_len) {
  if (str_len > INT_MAX) {
    return INT_MAX;
  }
  return str_len;
}

CQualityInfo toCQualityInfo(const eCAL::Registration::SQualityTopicInfo &quality
) {
  return {
      {quality.info.name.c_str(),
       quality.info.encoding.c_str(),
       quality.info.descriptor.c_str(),
       safe_len(quality.info.descriptor.size())},
      static_cast<uint8_t>(quality.quality),
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
