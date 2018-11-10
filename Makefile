PROJECT_NAME := "SBWeb"
PKG := "bmstu.codes/developers34/SBWeb"
PKG_LIST := $(go list ${PKG}/... | grep -v /vendor/)