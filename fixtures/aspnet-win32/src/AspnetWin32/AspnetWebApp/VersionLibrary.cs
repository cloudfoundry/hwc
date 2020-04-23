using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Diagnostics;
using System.IO;
using System.Linq;
using System.Runtime.InteropServices;
using System.Text;
using System.Web;

namespace NativeDll
{

    public static class Logger
    {
        public static void Info(string message)
        {
            Console.WriteLine($"-----> {message}");
        }
    }

    /// <summary>
    /// A peek into Win32 DLLs
    /// </summary>
    internal static class NativeLibrary
    {
        [DllImport("kernel32.dll", SetLastError = true)]
        [return: MarshalAs(UnmanagedType.Bool)]
        internal static extern bool FreeLibrary(IntPtr hModule);

#pragma warning disable CA2101 // Specify marshaling for P/Invoke string arguments
        [DllImport("kernel32", CharSet = CharSet.Ansi, ExactSpelling = false, SetLastError = true)]
#pragma warning restore CA2101 // Specify marshaling for P/Invoke string arguments
        internal static extern IntPtr LoadLibrary(string lpFileName);

#pragma warning disable CA2101 // Specify marshaling for P/Invoke string arguments
        [DllImport("kernel32", CharSet = CharSet.Ansi, ExactSpelling = true, SetLastError = true)]
#pragma warning restore CA2101 // Specify marshaling for P/Invoke string arguments
        internal static extern IntPtr GetProcAddress(IntPtr hModule, string procName);

        public static string GetLibraryPath(string dllFilename) => GetLibraryPath(dllFilename, System.Environment.Is64BitProcess);

        public static string GetLibraryPath(string dllFilename, bool is64bit)
        {
            string prefix = new[] { "Win32", "x64" }[Convert.ToInt16(is64bit)];

            var dllPath = System.IO.Path.GetFullPath(
                HttpContext.Current.Server.MapPath($@"bin\lib\{prefix}\{dllFilename}"));

            Logger.Info($"Test File.Exists({dllPath}) resulted in {File.Exists(dllPath)}");

            return dllPath;
        }
    }

    public class VersionLibraryWrapper : IDisposable
    {
        const string VersionDllFilename = "VersionDll.dll";
        public static VersionLibraryWrapper Create()
        {
            // Get 32-bit or 64-bit library directory
            var libPath = NativeLibrary.GetLibraryPath(VersionDllFilename);

            return new VersionLibraryWrapper(libPath);
        }

        public static VersionLibraryWrapper Createx86()
        {
            // Get 32-bit or 64-bit library directory
            var libPath = NativeLibrary.GetLibraryPath(VersionDllFilename, false);

            return new VersionLibraryWrapper(libPath);
        }

        public static VersionLibraryWrapper Createx64()
        {
            // Get 32-bit or 64-bit library directory
            var libPath = NativeLibrary.GetLibraryPath(VersionDllFilename, true);

            return new VersionLibraryWrapper(libPath);
        }

        // Delegate with function signature for the GetVersion function 
        [UnmanagedFunctionPointer(CallingConvention.Cdecl)]
        [return: MarshalAs(UnmanagedType.U4)]
        delegate UInt32 GetVersionDelegate(
            [OutAttribute][InAttribute] StringBuilder versionString,
            [OutAttribute] UInt32 length);

        // Handles and delegates
        IntPtr _dllhandle = IntPtr.Zero;
        GetVersionDelegate _getversion = null;

        public string GetVersion()
        {
            if (_getversion != null)
            {
                // Allocate buffer
                var size = 100;
                StringBuilder builder = new StringBuilder(size);

                // Get version string
                _getversion(builder, (uint)size);

                // Return string
                return builder.ToString();
            }

            return "";
        }

        public VersionLibraryWrapper(string libPath)
        {
            _dllhandle = NativeLibrary.LoadLibrary(libPath);

            // Handle error loading
            if (IntPtr.Zero == _dllhandle)
            {
                // Get the last error and display it.
                int win32error = Marshal.GetLastWin32Error();
                Logger.Info($"Unable to get handle on {libPath}. Win32 error {win32error}");
                return;
            }

            DelegateGetPlatformMessage();
        }

        private void DelegateGetPlatformMessage()
        {
            // Get handle to method in DLL
            var get_version_handle = NativeLibrary.GetProcAddress(_dllhandle, "GetPlatformMessage");

            // If successful, load function pointer
            if (get_version_handle != IntPtr.Zero)
            {
                _getversion = (GetVersionDelegate)Marshal.GetDelegateForFunctionPointer(
                    get_version_handle,
                    typeof(GetVersionDelegate));
            }
            else
            {
                int win32error = Marshal.GetLastWin32Error();
                Logger.Info($"Unable to get handle on GetProcAddress for GetPlatformMessage. Win32 error {win32error}");
            }
        }

        #region IDisposable Support
        private bool _disposedValue = false; // To detect redundant calls

        protected virtual void Dispose(bool disposing)
        {
            if (!_disposedValue)
            {
                if (disposing)
                {
                    // TODO: dispose managed state (managed objects).
                }

                // free unmanaged resources (unmanaged objects) 
                // ...set large fields to null.
                NativeLibrary.FreeLibrary(_dllhandle);

                _disposedValue = true;
            }
        }

        ~VersionLibraryWrapper()
        {
            // Do not change this code. Put cleanup code in Dispose(bool disposing) above.
            Dispose(false);
        }

        // This code added to correctly implement the disposable pattern.
        public void Dispose()
        {
            // Do not change this code. Put cleanup code in Dispose(bool disposing) above.
            Dispose(true);
            GC.SuppressFinalize(this);
        }
        #endregion
    }
}