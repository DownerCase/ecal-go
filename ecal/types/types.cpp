#include <climits>

#include "types.hpp"

namespace {
int safe_len(size_t str_len) {
  if (str_len > size_t{INT_MAX}) {
    return INT_MAX;
  }
  return static_cast<int>(str_len);
}
} // namespace

CDatatype toCType(const eCAL::SDataTypeInformation &datatype) {
  return {
      datatype.name.c_str(),
      datatype.encoding.c_str(),
      datatype.descriptor.data(),
      safe_len(datatype.descriptor.size()),
      {}
  };
}

CTopicId toCType(const eCAL::Registration::STopicId &id) {
  return {
      {
          id.topic_id.entity_id.data(),
          id.topic_id.host_name.data(),
          id.topic_id.process_id,
          {},
      },
      id.topic_name.data()
  };
}

CTopicMon toCType(const eCAL::Monitoring::STopicMon &topic) {
  return {
      topic.uname.c_str(),
      topic.tid.c_str(),
      topic.tname.c_str(),
      topic.direction.c_str(),
      toCType(topic.tdatatype),
      topic.dclock,
      topic.dfreq,
      topic.rclock,
      topic.tsize,
      topic.connections_loc,
      topic.connections_ext,
      topic.message_drops,
      {}
  };
}

CProcessMon toCType(const eCAL::Monitoring::SProcessMon &proc) {
  return {
      proc.hname.c_str(),
      proc.hgname.c_str(),
      proc.rclock,
      proc.pid,
      proc.pname.c_str(),
      proc.uname.c_str(),
      proc.pparam.c_str(),
      proc.state_severity,
      proc.state_severity_level,
      proc.state_info.c_str(),
      proc.component_init_info.c_str(),
      proc.ecal_runtime_version.c_str()
  };
}

CLogMessage toCType(const eCAL::Logging::SLogMessage &log) {
  return {
      log.time,
      log.hname.c_str(),
      log.pname.c_str(),
      log.uname.c_str(),
      log.content.c_str(),
      log.pid,
      log.level
  };
}
