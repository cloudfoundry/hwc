// VersionLibrary.cpp : Defines the entry point for the DLL application.
#include "pch.h"

#include <string>
#include <cstdint>
#include "VersionLibrary.h"


template<int>
std::string GetPlatform();

template<>
std::string GetPlatform<4>() { return "32-bit"; }

template<>
std::string GetPlatform<8>() { return "64-bit"; }

// helper function just to hide clumsy syntax
inline std::string GetPlatform() { return GetPlatform<sizeof(size_t)>(); }

class VersionLibrary
{
public:
	std::string GetPlatformMessage()
	{
		return "Native VersionLibrary: " + GetPlatform();
	}
};

VERSIONLIBRARY_API uint32_t GetPlatformMessage(char* buffer, uint32_t length)
{
	VersionLibrary vlib;
	auto val = snprintf(buffer, length, vlib.GetPlatformMessage().c_str());

	return 0;
}
