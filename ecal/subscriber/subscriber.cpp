#include "subscriber.h"

#include <ecal/ecal_subscriber.h>

#include "internal/handle_map.hpp"

extern "C" {
extern void goReceiveCallback(uintptr_t, void *, long);
}

namespace {
handle_map<eCAL::CSubscriber> subscribers;

void receive_callback(
    const uintptr_t handle,
    const eCAL::Registration::STopicId &topic,
    const eCAL::SDataTypeInformation &datatype,
    const eCAL::SReceiveCallbackData &data
) {
  goReceiveCallback(handle, data.buf, data.size);
}

} // namespace

bool NewSubscriber(uintptr_t handle) {
  const auto [it, added] = subscribers.emplace(handle);
  return added;
}

bool DestroySubscriber(uintptr_t handle) { return subscribers.erase(handle); }

bool SubscriberCreate(
    uintptr_t handle,
    const char *const topic,
    size_t topic_len,
    const char *const datatype_name,
    size_t datatype_name_len,
    const char *const datatype_encoding,
    size_t datatype_encoding_len,
    const char *const datatype_descriptor,
    size_t datatype_descriptor_len
) {
  auto *subscriber = subscribers.find(handle);
  if (subscriber == nullptr) {
    return false;
  }
  const auto created = subscriber->Create(
      std::string(topic, topic_len),
      {std::string(datatype_name, datatype_name_len),
       std::string(datatype_encoding, datatype_encoding_len),
       std::string(datatype_descriptor, datatype_descriptor_len)}
  );
  const auto bound_callback = [handle](
                                  const eCAL::Registration::STopicId &topic,
                                  const eCAL::SDataTypeInformation &datatype,
                                  const eCAL::SReceiveCallbackData &data
                              ) {
    receive_callback(handle, topic, datatype, data);
  };
  subscriber->AddReceiveCallback(bound_callback);
  return created;
}
