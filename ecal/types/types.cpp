#include "types.hpp"

namespace {
int safe_len(size_t str_len) {
  if (str_len > size_t{INT_MAX}) {
    return INT_MAX;
  }
  return static_cast<int>(str_len);
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
      {
          id.topic_id.entity_id.data(),
          id.topic_id.host_name.data(),
          id.topic_id.process_id,
      },
      id.topic_name.data()
  };
}
