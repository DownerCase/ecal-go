#include "types.hpp"

namespace {
int safe_len(size_t str_len) {
  if (str_len > INT_MAX) {
    return INT_MAX;
  }
  return str_len;
}
} // namespace

CDatatype toCDataType(const eCAL::SDataTypeInformation &datatype) {
  return {
      datatype.name.c_str(),
      datatype.encoding.c_str(),
      datatype.descriptor.data(),
      safe_len(datatype.descriptor.size())
  };
}

CTopicId toCTopicId(const eCAL::Registration::STopicId &id) {
  return {
      {id.topic_id.entity_id.data(),
       safe_len(id.topic_id.entity_id.size()),
       id.topic_id.process_id,
       id.topic_id.host_name.data(),
       safe_len(id.topic_id.host_name.size())},
      id.topic_name.data(),
      safe_len(id.topic_name.size())
  };
}
