#include <iostream>
#include <fstream>
#include <string>
#include <map>
#include <vector>
#include <cstdlib>

#ifdef _WIN32
#include <windows.h>
#include <process.h>
#else
#include <unistd.h>
#endif

//read environment variables from a file
std::map<std::string, std::string> readEnvVarsFromFile(const std::string& filePath);

//  set environment variables in the current process
void setEnvironmentVariables(const std::map<std::string, std::string>& envVars);

//execute the target program
bool executeTargetProgram(const std::string& programPath, const std::vector<std::string>& args);

int main(int argc, char* argv[]) {
    if (argc < 3) {
        std::cerr << "Usage: " << argv[0] << " <env_file_path> <target_executable> [args...]" << std::endl;
        return 1;
    }

    std::string envFilePath = argv[1];
    std::string targetProgram = argv[2];
    
    std::vector<std::string> programArgs;
    for (int i = 2; i < argc; i++) {
        programArgs.push_back(argv[i]);
    }

    auto envVars = readEnvVarsFromFile(envFilePath);
    
    setEnvironmentVariables(envVars);
    
    if (!executeTargetProgram(targetProgram, programArgs)) {
        std::cerr << "Failed to execute target program: " << targetProgram << std::endl;
        return 1;
    }
    
    return 0;
}