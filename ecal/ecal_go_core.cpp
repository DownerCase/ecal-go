#include "ecal_go_core.h"

#include <ecal/config/configuration.h>
#include <ecal/ecal_core.h>

namespace {
eCAL::Configuration convertConfig(config &) { return eCAL::Configuration{}; }
} // namespace

const char *GetVersionString() { return eCAL::GetVersionString(); }

const char *GetVersionDateString() { return eCAL::GetVersionDateString(); }

version GetVersion() {
  version version_{};
  // TODO: Version that uses refs instead of pointers
  eCAL::GetVersion(&version_.major, &version_.minor, &version_.patch);
  return version_;
}

int Initialize(struct config *config, const char *unit_name,
               unsigned int components) {
  auto cfg = convertConfig(*config);
  // TODO: Initialize should take by const ref
  return eCAL::Initialize(cfg, unit_name, components);
}

int Finalize() { return eCAL::Finalize(); }

bool IsInitialized(unsigned int component) {
  return eCAL::IsInitialized(component) == 1;
}

bool SetUnitName(const char *unit_name) {
  return eCAL::SetUnitName(unit_name) == 0;
}

bool Ok() { return eCAL::Ok(); }
