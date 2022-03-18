# 零售商注册账号
docker exec cli peer chaincode invoke -C vmichannel -n vmicc -c '{"Args":["retailerRegistration","lingshou1","5","3","20","8","2","9","25.9","12","5"]}'
# 供应商同意零售商注册
docker exec cli peer chaincode invoke -C vmichannel -n vmicc -c '{"Args":["supplierAuditRegistration","supplierAdmin","lingshou1","1"]}'
# 零售商查看供应商补货方案
docker exec cli peer chaincode invoke -C vmichannel -n vmicc -c '{"Args":["retailerViewScheme","lingshou1"]}'
# 零售商回应补货方案
docker exec cli peer chaincode invoke -C vmichannel -n vmicc -c '{"Args":["retailerResponseScheme","lingshou1","1"]}'
# 零售商更新库存
docker exec cli peer chaincode invoke -C vmichannel -n vmicc -c '{"Args":["retailerUpdateInventory","lingshou1","51"]}'
# 供货商查看零售商们补货方案
docker exec cli peer chaincode invoke -C vmichannel -n vmicc -c '{"Args":["supplierViewSchemes","supplierAdmin"]}'

docker exec cli peer chaincode invoke -C vmichannel -n vmicc -c '{"Args":["retailerRegistration","lingshou2","7","2","30","6","2","9","25.9","12","5"]}'
docker exec cli peer chaincode invoke -C vmichannel -n vmicc -c '{"Args":["supplierAuditRegistration","supplierAdmin","lingshou2","1"]}'
docker exec cli peer chaincode invoke -C vmichannel -n vmicc -c '{"Args":["retailerViewScheme","lingshou2"]}'

docker exec cli peer chaincode invoke -C vmichannel -n vmicc -c '{"Args":["retailerRegistration","lingshou3","7","2","30","6","2","9","25.9","12","5"]}'
docker exec cli peer chaincode invoke -C vmichannel -n vmicc -c '{"Args":["supplierAuditRegistration","supplierAdmin","lingshou3","0"]}'
docker exec cli peer chaincode invoke -C vmichannel -n vmicc -c '{"Args":["retailerViewScheme","lingshou3"]}'
docker exec cli peer chaincode invoke -C vmichannel -n vmicc -c '{"Args":["retailerViewScheme","lingshou4"]}'
