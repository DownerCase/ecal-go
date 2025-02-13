#include "publisher.h"

#include <ecal/pubsub/publisher.h>

void DestroyPublisher(void *publisher) {
  if (publisher == nullptr) {
    return;
  }
  delete static_cast<eCAL::CPublisher*>(publisher);
}

void *NewPublisher(
    const char *const topic,
    size_t topic_len,
    const char *const datatype_name,
    size_t datatype_name_len,
    const char *const datatype_encoding,
    size_t datatype_encoding_len,
    const char *const datatype_descriptor,
    size_t datatype_descriptor_len
) {
  return new eCAL::CPublisher(
      std::string(topic, topic_len),
      eCAL::SDataTypeInformation{
          std::string(datatype_name, datatype_name_len),
          std::string(datatype_encoding, datatype_encoding_len),
          std::string(datatype_descriptor, datatype_descriptor_len)
      }
  );
}

void PublisherSend(void *publisher, void *buf, size_t len) {
  if (publisher == nullptr) {
    return;
  }
  static_cast<eCAL::CPublisher*>(publisher)->Send(buf, len);
}
