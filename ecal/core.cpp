#include "core.h"

#include <ecal/config/configuration.h>
#include <ecal/core.h>
#include <ecal/defs.h>
#include <ecal/process.h>

extern "C" {
extern void goCopyString(uintptr_t, const char *);
}

const char *GetVersionString() { return ECAL_VERSION; }

const char *GetVersionDateString() { return ECAL_DATE; }

version GetVersion() {
  const auto version = eCAL::GetVersion();
  return {version.major, version.minor, version.patch};
}

bool Initialize(
    void *config,
    const char *unit_name,
    unsigned int components
) {
  auto* cfg = reinterpret_cast<eCAL::Configuration*>(config);
  return eCAL::Initialize(*cfg, unit_name, components);
}

bool Finalize() { return eCAL::Finalize(); }

bool IsInitialized() { return eCAL::IsInitialized(); }

bool IsComponentInitialized(unsigned int component) {
  return eCAL::IsInitialized(component);
}

bool Ok() { return eCAL::Ok(); }

void GetConfig(uintptr_t handle) {
  std::string cfg{};
  eCAL::Process::DumpConfig(cfg);
  goCopyString(handle, cfg.c_str());
}
