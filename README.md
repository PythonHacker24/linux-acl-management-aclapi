<div align="center">

# File Server Daemons Component (API Component) - Linux ACL Management Interface  

<img width="600" hegith="600" src="https://github.com/user-attachments/assets/a1625f58-0cd8-4df9-babc-31547b18d55a">

Securing Linux Storage with ACLs: An Open-Source Web Management Interface for Enhanced Data Protection.

A robust web-based management interface for Linux Access Control Lists (ACLs), designed to enhance data protection and simplify ACL administration. This project provides a modern, user-friendly solution for managing file system permissions in Linux environments.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

[View Documentation](https://pythonhacker24.github.io/linux-acl-management/)

</div>

## Project Summary 

Institutional departments, such as the Biomedical Informatics (BMI) Department of Emory University School of Medicine, manage vast amounts of data, often reaching petabyte scales across multiple Linux-based storage servers. Researchers storing data in these systems need a streamlined way to modify ACLs to grant or revoke access for collaborators. Currently, the IT team at BMI is responsible for manually handling these ACL modifications, which is time-consuming, error-prone, and inefficient, especially as data volume and user demands grow. To address this challenge at BMI and similar institutions worldwide, a Web Management Interface is needed to allow users to modify ACLs securely. This solution would eliminate the burden on IT teams by enabling on-demand permission management while ensuring security and reliability. The proposed system will feature a robust and highly configurable backend, high-speed databases, orchestration daemons for file storage servers, and an intuitive frontend. The proposal includes an in-depth analysis of required components, high-level and low-level design considerations, technology selection, and the demonstration of a functional prototype as proof of concept. The goal is to deliver a production-ready, secure, scalable, and reliable system for managing ACLs across multiple servers hosting filesystems such as NFS, BeeGFS, and others. This solution will streamline access control management and prepare it for deployment at BMI and other institutions worldwide, significantly reducing the manual workload for IT teams.

## Features

- Intuitive web interface for ACL management
- High-performance backend written in Go
- Real-time ACL updates
- Comprehensive ACL reporting and visualization
- Integration with OpenLDAP for authentication

## Development

### Branches

- `main`: Production-ready code
- `development-v<version>`: Development branches for specific versions

### ACL API Daemon 

The ACL API Daemom, a service called `aclapi` handles the gRPC connections with the backend component. 

It performs 2 jobs: 
1. Communicate with `aclcore` daemon about demanded ACL operations
2. Handler gRPC connections from backed component

It is provided with the least user privileges since it's exposed to the network.

Hence, this is not an independent component and needs `aclcore` to be running on the same server with proper setup.

Refer to the documentation for more information.

### Production Build (Manual)

For production build, it is recommended to use the Makefile. This allows you to build the complete binary on locally for security purposes. Since the project is in development mode, complete local build is not possible since dependencies are managed via GitHub and external vendors. Tarball based complete local builds will be developed in later stages.

Manual build provides more indepth look into how components are deployed and working. For automated deployment, use `install.sh` script.

1. Clone the repository:
   ```bash
   git clone https://github.com/PythonHacker24/linux-acl-management-aclapi.git
   cd linux-acl-management-aclapi
   ```

2. Use make:
    ```bash
    make build
    ```

3. Move the binary to /usr/local/bin
    ```bash
    sudo cp ./bin/aclapi /usr/local/bin/
    ```

4. Move configuration file to /etc/laclm
    ```bash
    sudo cp aclapi.yaml /etc/laclm/aclapi.yaml
    ```

5. Ensure you have a group called `laclm`. It should be followed in steps [aclcore daemon installation](https://github.com/PythonHacker24/linux-acl-management-aclcore).

6. Create a user called `aclapi` with no home, least privileges, and added into `laclm` group.

    ```bash
    sudo useradd --system --no-create-home --shell /usr/sbin/nologin --groups laclm aclapi
    ```

7. Create service for ACL API Daemon

    a. Create the systemd service file

    ```bash
    touch /etc/systemd/system/aclapi.service
    ```

    b. Copy this into aclapi.service

    ```ini
    [Unit]                                           
    Description=ACL API Daemon
    After=network.target laclm-daemon.service
    Requires=aclcore.service

    [Service]
    Type=simple

    ExecStart=/usr/local/bin/aclapi --config /etc/laclm/aclapi.yaml

    User=aclapi
    Group=laclm

    NoNewPrivileges=yes
    ProtectSystem=strict
    ProtectHome=yes
    PrivateTmp=yes

    Restart=on-failure

    [Install]
    WantedBy=multi-user.target
    ```

8. Reload SystemD daemons
    
    ```bash
    sudo systemctl daemon-reload
    ```

9. Enable aclapi service (optional: daemon will auto start when system is restarted)
    
    ```bash
    sudo systemctl enable aclapi.service
    ```

10. Start aclapi service
   
    ```bash
    sudo systemctl start aclapi.service
    ```

11. Check aclapi status 
    ```bash
    sudo systemctl status aclapi.service
    ```

## Project Structure

```
.
├── cmd/          # Application entry points
├── internal/     # Private application code
├── pkg/          # Public library code
├── api/          # API definitions and handlers
├── docs/         # Documentation
└── deployments/  # Deployment configurations
```

### Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct and development process.

## About

This project is developed as part of Google Summer of Code 2025, in collaboration with the Department of Biomedical Informatics at Emory University.

### Team

- **Contributor:** Aditya Patil
- **Mentors:** 
  - Robert Tweedy
  - Mahmoud Zeydabadinezhad, PhD

### Technologies

- **Backend:** Golang, net/http
- **API:** gRPC, REST
- **Infrastructure:** 
- **Packaging:** Tarball

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Department of Biomedical Informatics, Emory University
- Google Summer of Code Program
- Open Source Community

