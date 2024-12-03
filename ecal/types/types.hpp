#ifndef ECAL_GO_TYPES_HPP
#define ECAL_GO_TYPES_HPP

#include <vector>

#include <ecal/ecal_types.h>
#include <ecal/types/monitoring.h>

#include "types.h"

CDatatype toCType(const eCAL::SDataTypeInformation &datatype);
CTopicId toCType(const eCAL::Registration::STopicId &id);
CTopicMon toCType(const eCAL::Monitoring::STopicMon &topic);
CProcessMon toCType(const eCAL::Monitoring::SProcessMon &proc);

template <typename CType, typename EcalType, template<typename> typename Container >
std::vector<CType> toCTypes(const Container<EcalType> &ecaltypes) {
  std::vector<CType> ctypes{};
  ctypes.reserve(ecaltypes.size());
  for (const auto &ecaltype : ecaltypes) {
    ctypes.emplace_back(toCType(ecaltype));
  }
  return ctypes;
}


#endif
