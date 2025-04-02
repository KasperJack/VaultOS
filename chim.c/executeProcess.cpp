#include <windows.h>
#include <iostream>
#include <sstream>


char targetPath[MAX_PATH];


int main(int argc, char* argv[]) {


    GetModuleFileName(NULL, targetPath, MAX_PATH);

    char* lastSlash = strrchr(targetPath, '\\');
    if (lastSlash) *(lastSlash + 1) = '\0'; 

    strcat(targetPath, "sys.exe");

    std::cout << "Running: " << targetPath << std::endl;



    std::ostringstream cmdLineStream;
    cmdLineStream << targetPath;
    
    if (argc > 1) {
        for (int i = 1; i < argc; ++i) {
            cmdLineStream << " " << argv[i];
        }
    }
    
    std::string cmdLine = cmdLineStream.str();

    STARTUPINFO si = { sizeof(si) };
    PROCESS_INFORMATION pi;

    BOOL success = CreateProcess(
        NULL,                               
        const_cast<char*>(cmdLine.c_str()), 
        NULL,                               
        NULL,                              
        TRUE,                               
        0,                                  
        NULL,                               
        NULL,                               
        &si,                                
        &pi                                 
    );

    if (!success) {
        std::cerr << "Failed to launch CLI tool. Error code: " << GetLastError() << std::endl;
        return 1;
    }

    WaitForSingleObject(pi.hProcess, INFINITE);

    DWORD exitCode;
    GetExitCodeProcess(pi.hProcess, &exitCode);

    CloseHandle(pi.hProcess);
    CloseHandle(pi.hThread);

    return exitCode;
}
