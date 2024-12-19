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
    const eCAL::Registration::STopicId & /*topic*/,
    const eCAL::SDataTypeInformation & /*datatype*/,
    const eCAL::SReceiveCallbackData &data
) {
  goReceiveCallback(handle, data.buf, data.size);
}

} // namespace

bool NewSubscriber(
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

  const auto [it, added] = subscribers.emplace(
      handle,
      std::string(topic, topic_len),
      eCAL::SDataTypeInformation{
          std::string(datatype_name, datatype_name_len),
          std::string(datatype_encoding, datatype_encoding_len),
          std::string(datatype_descriptor, datatype_descriptor_len)
      }
  );
  if (!added) {
    return false;
  }

  const auto bound_callback = [handle](
                                  const eCAL::Registration::STopicId &_topic,
                                  const eCAL::SDataTypeInformation &_datatype,
                                  const eCAL::SReceiveCallbackData &_data
                              ) {
    receive_callback(handle, _topic, _datatype, _data);
  };
  return (*it).second.SetReceiveCallback(bound_callback);
}

bool DestroySubscriber(uintptr_t handle) { return subscribers.erase(handle); }
