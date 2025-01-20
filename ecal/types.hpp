#ifndef ECAL_GO_TYPES_HPP
#define ECAL_GO_TYPES_HPP

#include <vector>

#include <ecal/ecal_types.h>
#include <ecal/types/monitoring.h>

#include "types.h"

CDatatype toCType(const eCAL::SDataTypeInformation &datatype);

template <
    template <typename> typename DstContainerT,
    typename DstT,
    template <typename> typename SrcContainerT,
    typename SrcT>
std::vector<DstT>
containerTo(const SrcContainerT<SrcT> &src, DstT (*proj)(const SrcT &)) {
  DstContainerT<DstT> dst{};
  dst.reserve(src.size());
  for (const auto &ecaltype : src) {
    dst.emplace_back(proj(ecaltype));
  }
  return dst;
}

template <
    typename CType,
    typename EcalType,
    template <typename> typename Container>
std::vector<CType> toCTypes(const Container<EcalType> &ecaltypes) {
  return containerTo<std::vector>(ecaltypes, toCType);
}

#endif
