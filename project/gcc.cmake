set(CMAKE_C_COMPILER gcc)
set(CMAKE_CXX_COMPILER g++)

set(gcc_warnings
-Wpedantic               # Forbid extensions
-pedantic-errors         # Warnings from Wpedantic (and default warnings) are errors
-Wall                    # Enable warnings
-Wextra                  # Enable more warnings
-Wdouble-promotion       # Accidental float -> double promotion
-Wformat=2               # Check arguments to printf etc.
-Wformat-overflow=2      # Check for potential sprintf etc. destination overflows
-Wformat-signedness      # Checking format string signedness
-Wformat-truncation=2    # Check for output truncation
-Wnull-dereference       # Dereference of null
-Wnrvo                   # Warn when NRVO fails
-Wimplicit-fallthrough=5 # Fallthrough can only be allowed by [[fallthrough]]
-Wmissing-include-dirs   # Include directories that don't exist
-Wswitch-default         # Require default for switch
-Wswitch-enum            # Require all named cases for switch on an enum
-Wunused
-Wuse-after-free=3
-Wuseless-cast           # Casting to your own type
-Wmaybe-uninitialized    # Potential unitialized uses
-Wstrict-overflow=5
-Wstringop-overflow=4
-Warith-conversion       # WARN: Noisy warning
-Warray-bounds=2
-Wbidi-chars=any         # Misleading bidirectional UTF-8 control characters
-Wduplicated-branches    # Duplicated conditional branches
-Wduplicated-cond        # Duplicated conditions
-Wtrampolines            # No trampolines
-Wfloat-equal            # Comparing floats is hard
-Wshadow                 # Don't shadow varibles
-Wunsafe-loop-optimizations
-Wundef
-Wcast-qual
-Wcast-align=strict
-Wconversion             # Catch implicit conversions
-Wdate-time              # Date/time prevents reproducible builds
-Wsign-conversion        # Implicit conversions that change the sign
-Wlogical-op             # Suspicious use of logical operators
-Wmissing-declarations
-Wpadded
-Wno-error=padded        # WARN: Wpadded is extremely noisy
-Wredundant-decls
-Winline
-Wdisabled-optimization
-Wstack-protector

-Wctad-maybe-unsupported
-Wctor-dtor-privacy
-Wnoexcept
-Wnon-virtual-dtor       # Virtual classes need virtual destructors
-Wredundant-tags
-Weffc++
-Wstrict-null-sentinel
-Wold-style-cast
-Woverloaded-virtual
-Wsign-promo
-Wcatch-value=3
-Wextra-semi

-Wsuggest-final-types
-Wsuggest-final-methods
-Wsuggest-override
)

string(REPLACE ";" " " gcc_warnings "${gcc_warnings}")

set(CMAKE_CXX_FLAGS_INIT "${gcc_warnings}")

set(CMAKE_COMPILE_WARNING_AS_ERROR ON)
