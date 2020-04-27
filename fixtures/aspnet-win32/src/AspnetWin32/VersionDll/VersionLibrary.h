// VersionLibrary.h - Contains declarations of math functions
#pragma once
#include <cstdint>

#ifdef VERSIONDLL_EXPORTS
#define VERSIONLIBRARY_API __declspec(dllexport)
#else
#define VERSIONLIBRARY_API __declspec(dllimport)
#endif // VERSIONDLL_EXPORTS

extern "C" VERSIONLIBRARY_API uint32_t GetPlatformMessage(char* buffer, uint32_t length);