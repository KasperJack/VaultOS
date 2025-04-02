#include <windows.h>
#include <iostream>
#include <fstream>
#include <unordered_map>
#include <string>

std::string getDriveLetter() {
    char path[MAX_PATH];
    if (GetModuleFileNameA(NULL, path, MAX_PATH)) {
        return std::string(1, path[0]);
    }
    return "";
}

std::string trim(const std::string& str) {
    size_t first = str.find_first_not_of(" \t\r\n");
    if (first == std::string::npos) return "";
    size_t last = str.find_last_not_of(" \t\r\n");
    return str.substr(first, (last - first + 1));
}

std::unordered_map<std::string, std::string> loadEnv(const std::string& filename, const std::string& driveLetter) {
    std::unordered_map<std::string, std::string> envVars;
    std::ifstream file(filename);

    if (!file) {
        std::cerr << "Failed to open env file: " << filename << std::endl;
        return envVars;
    }

    std::string line;
    while (std::getline(file, line)) {
        line = trim(line); // Remove leading/trailing spaces

        // Ignore comments and empty lines
        if (line.empty() || line[0] == '#') continue;

        // Find '=' to split key and value
        size_t pos = line.find('=');
        if (pos == std::string::npos) continue; // Skip invalid lines

        std::string key = trim(line.substr(0, pos));
        std::string value = trim(line.substr(pos + 1));

        if (value.front() == '"' && value.back() == '"') {
            value = value.substr(1, value.size() - 2);
        }

        //  $DRIVE_LETTER
        size_t drivePos;
        while ((drivePos = value.find("$DRIVE_LETTER")) != std::string::npos) {
            value.replace(drivePos, 13, driveLetter);
        }

        envVars[key] = value;
    }

    return envVars;
}

int main() {
    std::string driveLetter = getDriveLetter();
    if (driveLetter.empty()) {
        std::cerr << "Failed to determine drive letter!" << std::endl;
        return 1;
    }

    std::cout << "Detected Drive Letter: " << driveLetter << ":\\" << std::endl;

    std::unordered_map<std::string, std::string> envVars = loadEnv("config.env", driveLetter);

    std::cout << "\nLoaded Environment Variables:\n";
    for (const auto& [key, value] : envVars) {
        std::cout << key << " = " << value << std::endl;
    }

    return 0;
}
