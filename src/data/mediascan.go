package data

func (d *DataLayerPg) MediaScan(musicFolders []string, count chan<- int, done chan<- struct{}) {
	defer func() {
		done <- struct{}{}
	}()

	//TODO: implement me

	/*
		for each music folders
			scan for artists
			 - folder name
			 - cover art?
			 - how many subfolders? -> album count

			 for each artist
			 	scan for albums
				 - folder name
				 - cover art
				 - how many music files

				for each album
					scan for songs
					 - title
					 - is_dir
					 - cover_art
					 - bit_rate
					 - size
					 - suffix
					 - content_type
					 - is_video

		export to db
		export artist -> get id
			export album -> get id
				export songs
	*/
}
