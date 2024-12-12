#include <iostream>
#include <string>
#include <fstream>
#include <sstream>
#include <vector>
#include <filesystem>

std::string get_resource_path() {
#ifdef _WIN32
    const char* appdata = std::getenv("APPDATA");
    if (appdata)
        return std::string(appdata) + "\\LLWC\\resources";
    else
        std::cerr << "Could not find \%APPDATA\%\\LLWC\\resources directory" << std::endl;
#else
    const char* home = std::getenv("HOME");
    if (home)
        return std::string(home) + "/.LLWV/resources";
    else
        std::cerr << "Could not find ~/.LLWC/resources directory" << std::endl;
#endif
    return "resources";
}

std::vector<std::vector<std::string>> parseCSV(const std::string& filepath) {
    std::vector<std::vector<std::string>> data;

    if(!std::filesystem::exists(filepath)) {
        throw std::runtime_error("File not found: " + filepath);
    }

    std::fstream file(filepath);

    if (!file.is_open()) {
        throw std::runtime_error("Unable to open file: " + filepath);
    }

    std::string line;
    while(std::getline(file, line)) {
        std::vector<std::string> row;
        std::stringstream lineStream(line);
        std::string cell;

        while(std::getline(lineStream, cell, ',')) {
            row.push_back(cell);
        }

        data.push_back(row);
    }

    file.close();

    return data;
}


int main() {
    std::cout << "Loading resources from: " << get_resource_path() << std::endl;;

    try {
        const std::string filepath = "../resources/pokemon.csv";
        auto csvData = parseCSV(filepath);

        for(const auto& row : csvData) {
            for(const auto& cell : row) {
                std::cout << cell << " ";
            }
            std::cout << std::endl;
        }
    } catch(const std::exception& e) {
        std::cerr << "Error: " << e.what() << std::endl;
    }

    return 0;
}