#ifndef ECAL_GO_TYPES_HPP
#define ECAL_GO_TYPES_HPP

#include <climits>
#include <ecal/ecal_types.h>

#include "types.h"

CDatatype toCDataType(const eCAL::SDataTypeInformation &datatype);
CTopicId toCTopicId(const eCAL::Registration::STopicId &id);

#endif
