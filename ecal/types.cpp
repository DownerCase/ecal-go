#include <climits>
#include <ecal/ecal_process_severity.h>
#include <ecal/types/monitoring.h>

#include "process.h"
#include "types.h"
#include "types.hpp"

namespace {
int safe_len(size_t str_len) {
  if (str_len > size_t{INT_MAX}) {
    return INT_MAX;
  }
  return static_cast<int>(str_len);
}

} // namespace

CDatatype toCDatatype(const eCAL::SDataTypeInformation &datatype) {
  return {
      datatype.name.c_str(),
      datatype.encoding.c_str(),
      datatype.descriptor.data(),
      safe_len(datatype.descriptor.size()),
      {}
  };
}

static_assert(
    process_severity_unknown ==
    static_cast<CProcessSeverity>(eCAL::Process::eSeverity::unknown)
);
static_assert(
    process_severity_healthy ==
    static_cast<CProcessSeverity>(eCAL::Process::eSeverity::healthy)
);
static_assert(
    process_severity_warning ==
    static_cast<CProcessSeverity>(eCAL::Process::eSeverity::warning)
);
static_assert(
    process_severity_critical ==
    static_cast<CProcessSeverity>(eCAL::Process::eSeverity::critical)
);
static_assert(
    process_severity_failed ==
    static_cast<CProcessSeverity>(eCAL::Process::eSeverity::failed)
);
