package mysql

//model := bitcoinModel.NewConnection()
//
//model.StartTransaction().
//	StoreAddress(wallet.walletId, coinId, changeAddr, 0, true, false).
//	StoreAddress(wallet.walletId, coinId, addr,       0, false, false).
//	Commit()
//defer func(model *bitcoinModel.Connection){
//	err := model.Rollback()
//	if err!=nil { retErr = ErrorsDepth(2, err, ErrDbOpFail) }
//}(model)
//if model.Error()!=nil {
//	return nil
//}
