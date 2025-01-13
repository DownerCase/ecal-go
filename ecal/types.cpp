#include <climits>
#include <ecal/types/monitoring.h>

#include "types.h"
#include "types.hpp"

namespace {
int safe_len(size_t str_len) {
  if (str_len > size_t{INT_MAX}) {
    return INT_MAX;
  }
  return static_cast<int>(str_len);
}

template <class T> CServiceCommon toServiceCommon(const T &t) {
  return {
      t.sname.c_str(),
      t.sid,
      nullptr, // Methods are filled in a separate pass
      t.methods.size(),
      t.hname.c_str(),
      t.pname.c_str(),
      t.uname.c_str(),
      t.rclock,
      t.pid,
      t.version,
      {}
  };
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
          id.topic_id.entity_id,
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
      topic.hname.c_str(),
      topic.tid,
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

CClientMon toCType(const eCAL::Monitoring::SClientMon &client) {
  return {toServiceCommon(client)};
}

CServerMon toCType(const eCAL::Monitoring::SServerMon &server) {
  return {toServiceCommon(server), server.tcp_port_v0, server.tcp_port_v1};
}

CMethodMon toCType(const eCAL::Monitoring::SMethodMon &method) {
  return {
      method.mname.c_str(),
      toCType(method.req_datatype),
      toCType(method.resp_datatype),
      method.call_count
  };
}
