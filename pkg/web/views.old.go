package web

//func (v *views) getBook(c echo.Context) error {
//	bookID := c.Param("book")
//
//	book, err := v.game.GetBookByID(bookID)
//	if err != nil {
//		return errors.Wrap(err, "failed to get book")
//	}
//
//	path := fmt.Sprintf("/p/%s/start", book.ID)
//	return c.Redirect(http.StatusFound, path)
//}

//const initFlag = "--init-has-happened--"

//func (v *views) getPage(c echo.Context) error {
//	bookID := c.Param("book")
//	pageID := c.Param("page")
//	userID := getUserID(c)
//	playerStorage := storage.NamespacedStorage(v.storage, userID)
//
//	book, err := v.game.GetBookByID(bookID)
//	if err != nil {
//		return errors.Wrap(err, "failed to get book")
//	}
//
//	bookResults, err := v.player.ExecuteBook(book, playerStorage)
//	if err != nil {
//		return errors.Wrap(err, "failed to execute book")
//	}
//
//	if pageID == "start" {
//		path := fmt.Sprintf("/p/%s/%s", book, bookResults.StartPage)
//		return c.Redirect(http.StatusFound, path)
//	}
//
//	if item := playerStorage.Get(initFlag); item == nil {
//		if err = v.player.BeginBook(book, playerStorage); err != nil {
//			return errors.Wrap(err, "failed to init book")
//		}
//		playerStorage.Set(initFlag, true)
//	}
//
//	page, err := v.game.GetPage(bookID, pageID)
//	if err != nil {
//		return errors.Wrapf(err, "failed to get page bookID=%s/pageID=%s", bookID, pageID)
//	}
//
//	pageResults, err := v.player.ExecutePage(book, page, playerStorage)
//	if err != nil {
//		return errors.Wrapf(err, "failed to get page bookID=%s/pageID=%s", bookID, pageID)
//	}
//
//	viewModel, err := v.generatePageViewModel(pageResults, bookResults, page)
//	if err != nil {
//		return errors.Wrapf(err, "failed to generate viewModel bookID=%s/pageID=%s", bookID, pageID)
//	}
//
//	return v.renderTemplate(c, "page.gohtml", viewModel)
//}
