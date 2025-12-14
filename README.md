[![MusicStreaming CI](https://github.com/FGguy/MusicStreaming/actions/workflows/CI.yaml/badge.svg)](https://github.com/FGguy/MusicStreaming/actions/workflows/CI.yaml)

# MusicStreaming

_A self-hosted music streaming platform_

## About

MusicStreaming is an open-source, self-hosted solution for streaming your personal music collection. It’s designed so you can control your audio library, stream to your devices, and not rely on third-party cloud services.

---

## Features

- Uses the Subsonic API, use any Subsonic compatible client.
- Serve your local music files (MP3, FLAC, etc) over the network.
- Web-based interface to browse artists, albums, playlists.
- User authentication / multi-user support.
- Streaming via HTTP/HTTPS.
- Lightweight and easy to run on home servers or small VPS.

---

## Architecture / Tech Stack

- Monolithic architecture.
- Postgres for data persistance.
- Redis for caching data and reducing latency.
- Docker for containerization.
- Docker Compose for deployment.

---

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

## Contact

If you have questions or suggestions, feel free to open an issue or contact me (@FGguy) via GitHub.

---
