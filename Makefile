image:
	@cd jqtransformation && gcloud builds submit --tag gcr.io/ultra-hologram-297914/jqt
	@cd mongodbtarget && gcloud builds submit --tag gcr.io/ultra-hologram-297914/mongodbtarget
