#include <windows.h>
#include <iostream>
#include <fstream>
#include <unordered_map>
#include <string>
#include <regex>

std::string getDriveLetter() {
    char path[MAX_PATH];
    if (GetModuleFileNameA(NULL, path, MAX_PATH)) {
        return std::string(1, path[0]) + ":";
    }
    return "";
}

std::string trim(const std::string& str) {
    size_t first = str.find_first_not_of(" \t\r\n");
    if (first == std::string::npos) return "";
    size_t last = str.find_last_not_of(" \t\r\n");
    return str.substr(first, (last - first + 1));
}

std::string fixBackslashes(std::string path) {
    size_t pos = 0;
    while ((pos = path.find("\\\\", pos)) != std::string::npos) {
        path.replace(pos, 2, "\\");
    }
    return path;
}

std::string expandVariables(const std::string& value, const std::unordered_map<std::string, std::string>& envVars) {
    std::string expanded = value;
    std::regex varPattern(R"(\$([A-Za-z_]+))"); // Matches variables like $VAR_NAME
    std::smatch match;

    while (std::regex_search(expanded, match, varPattern)) {
        std::string varName = match[1].str(); // Extract variable name
        auto it = envVars.find(varName);
        if (it != envVars.end()) {
            expanded.replace(match.position(0), match.length(0), it->second); // Replace with value
        } else {
            expanded.replace(match.position(0), match.length(0), ""); // Remove undefined variables
        }
    }

    return expanded;
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

        if (line.empty() || line[0] == '#') continue; // Skip comments and empty lines

        size_t pos = line.find('=');
        if (pos == std::string::npos) continue; // Skip invalid lines

        std::string key = trim(line.substr(0, pos));
        std::string value = trim(line.substr(pos + 1));

        if (value.front() == '"' && value.back() == '"') {
            value = value.substr(1, value.size() - 2); // Remove quotes
        }

        size_t drivePos;
        while ((drivePos = value.find("$DRIVE_LETTER")) != std::string::npos) {
            value.replace(drivePos, 13, driveLetter);
        }

        envVars[key] = value; // Store in map
    }

    // Expand all stored variables in values
    for (auto& [key, val] : envVars) {
        val = expandVariables(val, envVars);
        // Fix backslashes
        val = fixBackslashes(val);
    }

    return envVars;
}

int main() {
    std::string driveLetter = getDriveLetter();
    if (driveLetter.empty()) {
        std::cerr << "Failed to determine drive letter!" << std::endl;
        return 1;
    }

    std::cout << "Detected Drive Letter: " << driveLetter << std::endl;

    std::unordered_map<std::string, std::string> envVars = loadEnv("config.env", driveLetter);

    std::cout << "\nLoaded Environment Variables:\n";
    for (const auto& [key, value] : envVars) {
        std::cout << key << " = " << value << std::endl;
    }

    return 0;
}
