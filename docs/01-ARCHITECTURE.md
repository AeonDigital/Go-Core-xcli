ARCHITECTURE
================================

1. OVERVIEW, PHILOSOPHY AND COMPLIANCE

This document establishes the structural, architectural, and documentation guidelines for Go modules within this ecosystem.
Our core development philosophy balances simplicity with rigid code boundaries. We value explicit design over implicit behaviors, ensuring zero black boxes. Every component must have a predictable location, clear constraints, and maximum human readability.

COMPLIANCE RULE:
This architecture is preferred and strongly encouraged for all projects to maintain ecosystem consistency. However, it is not mandatory. If a project is too small, or if this structure introduces accidental complexity and makes development extremely difficult, developers are allowed to adapt or simplify this layout.



&nbsp;
&nbsp;


________________________________________________________________________________

## 2. REPOSITORY STRUCTURE

Replace the pkgname root directory with your module name, and substitute the pkg placeholder with the principal shorthand name of your library to enforce context-prefixed public packages.

When a public package mirrors the core purpose of an internal one (e.g., pkgcfg vs internalconfig), the public layer must act exclusively as an exposed gateway or contract, handling data primitives or public structures, while delegating heavy orchestration to its private counterpart.

```
mainrepo
└── pkgname
    ├── docs                  
    │   ├── 00-MD_RULES.md      
    │   └── 01-ARCHITECTURE.md  
    │
    ├── internal                # PRIVATE: Hidden core, zero external imports allowed
    │   ├── config              # Heavy loaders (flags, envs, file parsers)
    │   ├── contt               # Internal immutable constants and private magic numbers
    │   ├── fn                  # Core stateless algorithms and computational math
    │   ├── intfc               # Private decoupled contracts
    │   ├── struc               # Domain models, strict validations, and stateful structures
    │   └── xerrors             # Private named sentinel errors
    │
    └── pkg                     # PUBLIC: Exportable gateway, exposes primitives and contracts
        └── pkgname          
            ├── pkgcfg          # Public configuration structures and primitives
            ├── pkgconstt       # Public type definitions, enums, and validation constants
            ├── pkgfn           # Public stateless utility wrappers (Print, Formatter)
            ├── pkgintfc        # Public behavioral interfaces for external extendability
            ├── pkgstruc        # Public DTOs and data models needed by external callers
            └── pkgerrors       # Publicly interceptable named sentinel errors
```


&nbsp;
&nbsp;


________________________________________________________________________________

## 3. COMPONENT DIRECTORY RULES

### /internal

Contains the private core of the application or library. The Go compiler strictly forbids any external micro-project or module from importing code residing here.

* internal/config/ (package config)
  Entry point for parsing flags, environment variables, or files. No business logic allowed.
* internal/contt/ (package contt)
  Central registry for immutable primitive values. Mutable global states or pointers are strictly banned.
* internal/fn/ (package fn)
  Hosts stateless, deterministic computational routines and mathematical algorithms.
* internal/intfc/ (package intfc)
  Defines decoupled behavioral contracts. Interfaces must be small (1 to 3 methods).
* internal/struc/ (package struc)
  Hosts pure data models and DTOs. Methods are limited to basic validation or getters/setters.
* internal/xerrors/ (package xerrors)
  Central registry for named sentinel errors (e.g., ErrPermissionDenied).


&nbsp;


### /pkg

Contains the public, exportable interface of the library or application. Other micro-projects within the monorepo can freely import packages from this directory.

* Architectural Symmetry Rule: The 6 core functional layers (config, constt, fn, intfc, struc, xerrors) can exist both in /internal (for private orchestration) and /pkg (for public exposure).
* Context-Prefixed Naming Constraint: To prevent package name collisions, any of the 6 core layers exposed in /pkg must be prefixed with the shorthand name of the project (e.g., <pkg>cfg, <pkg>constt, <pkg>fn, <pkg>intfc, <pkg>struc, <pkg>errors).
* Public Constants & Errors Rule: Layers like <pkg>constt and <pkg>errors are explicitly designed to expose types, sentinel errors, and validation primitives that external consumers or package objects need to interact with the library's API.
