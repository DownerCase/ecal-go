#include "core.h"

#include <ecal/config/configuration.h>
#include <ecal/ecal_core.h>
#include <ecal/ecal_defs.h>

namespace {
eCAL::Configuration convertConfig(config &) { return eCAL::Configuration{}; }
} // namespace

const char *GetVersionString() { return ECAL_VERSION; }

const char *GetVersionDateString() { return ECAL_DATE; }

version GetVersion() {
  const auto version = eCAL::GetVersion();
  return {version.major, version.minor, version.patch};
}

int Initialize(config *config, const char *unit_name, unsigned int components) {
  auto cfg = convertConfig(*config);
  // TODO: Initialize should take by const ref
  return eCAL::Initialize(cfg, unit_name, components);
}

int Finalize() { return eCAL::Finalize(); }

bool IsInitialized() { return eCAL::IsInitialized() == 1; }

bool IsComponentInitialized(unsigned int component) {
  return eCAL::IsInitialized(component) == 1;
}

bool SetUnitName(const char *unit_name) {
  return eCAL::SetUnitName(unit_name) == 0;
}

bool Ok() { return eCAL::Ok(); }
