#include "subscriber.h"

#include <ecal/ecal_subscriber.h>
#include <stdexcept>

#include "internal/handle_map.hpp"

namespace {
handle_map<eCAL::CSubscriber> subscribers;
handle_map<std::string> messages;
} // namespace

const void *NewSubscriber() {
  const auto [it, added] = subscribers.emplace();
  if (!added) {
    return nullptr;
  }
  return it->second.get();
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
  return subscriber->Create(
      std::string(topic, topic_len),
      {std::string(datatype_name, datatype_name_len),
       std::string(datatype_encoding, datatype_encoding_len),
       std::string(datatype_descriptor, datatype_descriptor_len)}
  );
}

uintptr_t
SubscriberReceive(uintptr_t subscriber_handle, const char **msg, size_t *len) {
  auto *subscriber = subscribers.find(subscriber_handle);
  if (subscriber == nullptr) {
    *msg = nullptr;
    *len = 0;
    return 0;
  }

  // Receive message into our buffer
  // TODO: Replace with callback based method to remove a copy
  std::string buffer{};
  const auto received = subscriber->ReceiveBuffer(buffer, nullptr, -1);
  if (!received) {
    *msg = nullptr;
    *len = 0;
    return 0;
  }

  // Save the message for later processing
  auto [it, added] = messages.emplace(std::move(buffer));
  if (!added) {
    throw std::runtime_error("Failed to store received message");
    *msg = nullptr;
    *len = 0;
    return 0;
  }

  *msg = it->second->c_str();
  *len = it->second->size();

  return it->first;
}

bool DestroyMessage(uintptr_t message_handle) {
  return messages.erase(message_handle);
}
