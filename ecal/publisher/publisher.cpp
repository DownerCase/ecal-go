#include "publisher.h"

#include <ecal/pubsub/publisher.h>

#include "internal/handle_map.hpp"

namespace {
handle_map<eCAL::CPublisher> publishers{};
} // namespace

bool DestroyPublisher(uintptr_t handle) { return publishers.erase(handle); }

bool NewPublisher(
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
  const auto [it, added] = publishers.emplace(
      handle,
      std::string(topic, topic_len),
      eCAL::SDataTypeInformation{
          std::string(datatype_name, datatype_name_len),
          std::string(datatype_encoding, datatype_encoding_len),
          std::string(datatype_descriptor, datatype_descriptor_len)
      }
  );
  return added;
}

void PublisherSend(uintptr_t handle, void *buf, size_t len) {
  auto *publisher = publishers.find(handle);
  if (publisher == nullptr) {
    return;
  }
  publisher->Send(buf, len);
}
