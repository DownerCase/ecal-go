#include "registration.h"

#include "types.hpp"

#include <climits>
#include <ecal/ecal_registration.h>

extern "C" {
extern void
goPublisherEventCallback(uintptr_t handle, struct CTopicId id, uint8_t event);
extern void goPublisherDatatypeCallback(
    uintptr_t handle,
    struct CTopicId id,
    uint8_t event,
    struct CQualityInfo quality
);
}

namespace {
int safe_len(size_t str_len) {
  if (str_len > INT_MAX) {
    return INT_MAX;
  }
  return str_len;
}

struct CQualityInfo
toCQualityInfo(const eCAL::Registration::SQualityTopicInfo &quality) {
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
  const auto callback_adapter =
      [handle](
          const eCAL::Registration::STopicId &id,
          eCAL::Registration::RegistrationEventType event
      ) {
        goPublisherEventCallback(
            handle,
            toCTopicId(id),
            static_cast<uint8_t>(event)
        );
      };
  return eCAL::Registration::AddPublisherEventCallback(callback_adapter);
}

void RemPublisherEventCallback(size_t handle) {
  eCAL::Registration::RemPublisherEventCallback(handle);
}
