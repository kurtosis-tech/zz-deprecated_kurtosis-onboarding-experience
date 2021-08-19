Ethereum On-Boarding Testsuite
==============================

## Implement a Single Node Ethereum Test Network

1. Create an account on https://kurtosistech.com if you don't have one yet.
1. Verify that the Docker daemon is running on your local machine.
1. Clone this repository by running `git clone https://github.com/kurtosis-tech/kurtosis-onboarding-experience.git`
1. Run the empty Ethereum single node test `my_test` to verify that everything works on your local machine.
   1. Run `bash scripts/build-and-run.sh all`
   1. Verify that the output of the build-and-run.sh script indicates that one test ran (my_test) and that it passed.
1. Import the dependencies that are used in this example test suite.
   1. Run `go get github.com/ethereum/go-ethereum`
   1. Run `go get github.com/palantir/stacktrace`
   1. Run `go get github.com/sirupsen/logrus`
1. Write the Setup() method on the Ethereum single node test in order to bootstrap the network and leave it running and ready to use it.
   1. In your preferred IDE, open the Ethereum single node test `my_test` at `testsuite/testsuite_impl/my_test/my_test.go`
   1. Implement the private helper function `getEthereumServiceConfigurations` used to set the Ethereum node service inside the testsuite
      1. Add the following private helper functions `getContainerCreationConfig()`, `getRunConfigFunc()` and `getEthereumServiceConfigurations()` in [this Gist](https://gist.github.com/leoporoli/123cb1d682d74dcafe7920f01809b167) to the bottom of the file, so they can be used later. 
   1. Implement the first part of the Setup() method which starts the Ethereum single node service in the testsuite
      1. Add the following code lines in [this Gist](https://gist.github.com/leoporoli/d81577dcfd5fdc6605ccdf9f61eed81f) in the top of the Setup() method
   1. Implement the second and last part of the Setup() method to check if the service is available
      1. Add the following private helper functions `getEnodeAddress()` and `sendRpcCall` in [this Gist](https://gist.github.com/leoporoli/1a03e9500a61a20d06ed8e3827d72f5e) to the bottom of the file, so they can be used later.
      1. Add the following code lines in [this Gist](https://gist.github.com/leoporoli/f9aacad32b2800a98f68bcf5fa32165c) in the bottom of the Setup() method
   1. Verify that running `bash scripts/build-and-run.sh all` generates output indicating that one test ran (my_test) and that it passed
1. Write the Run() method on the Ethereum single node test in order to test the initilization of an Ethereum Single Node Network in Dev Mode
   1. Read [the official Geth documentation for start Ethereum in Dev mode](https://geth.ethereum.org/docs/getting-started/dev-mode) to have a context as this test implements the commands that are described in this document
   1. Implement the first part of the Run() method which get the Ethereum single node service from the network, get the Ethereum client from the service and uses this to get the chain ID of the private Ethereum Network
      1. Add the private helper function `getEthClient` in [this Gist](https://gist.github.com/leoporoli/8a1641f9d78752f984ed672895c1f97c) to the bottom of the file, so it can be used later.
      1. Add the following code lines in [this Gist](https://gist.github.com/leoporoli/628c78452e7e05549809f5cbcb62cdf6) in the top of the Run() method
      1. Verify that running `bash scripts/build-and-run.sh all` generates output indicating that one test ran (my_test) and that it passed
   1. Implement the second and last part of the Run() method which implements all the commands of the official Geth documentation in the `Dev mode` section
      1. 1. Add the following code lines in [this Gist](https://gist.github.com/leoporoli/5a1539d2e8e45d3658a5b2398d9f3ba7) in the bottom of the Run() method
   1. Verify that running `bash scripts/build-and-run.sh all` generates output indicating that one test ran (my_test) and that the test contained business logic for an Ethereum single node network test, and that it passed.

## Implement an Advanced Test which test and Ethereum Private Network with Multiple Nodes

El objetivo de este test lo podemos dividir en dos partes:

La primera parte es testear una implementación de una red privada de Ethereum con multiples nodos, que utilice `Clique consensus` como prueba de autoridad y que este previamente seteado en el genesis block, con una cuenta firmante.
Para iniciar la red privada primero será necesario iniciar el bootnode y luego iniciar los nodos restantes utilizando el bootnode, al final evaluar la cantidad de peers que contiene cada nodo para validar que los nodos de la red están correctamente interconectados
Para esta primera parte se puede apoyar en la documentación oficial sobre [como setear una red privada](https://geth.ethereum.org/docs/interface/private-network) provista por los desarrolladores de Geth

La segunda parte del test consiste en testear ejecuciones dentro de la red privada previamente seteada. Para esto primero se deberá deployar el smart contract `Hello World` que se encuentra en la carpeta `smart_contracts/bindings/hello_world.go` y validar que el mismo ha sido exitosamente deployado.
Y por último validar que la cuenta firmante se encuentra dentro de la red privada.

1. Implementar el método Setup() del test `my_advanced_test_.go` para inicializar la red Ethereum privada con múltiples nodos
   1. In your preferred IDE, open the advanced Ethereum test `my_advanced_test` at `testsuite/testsuite_impl/my_advanced_test/my_advanced_test.go`
   1. Inicie la red privada que contenga el genesis block que se encuentra dentro de `testsuite/data/genesis.json`, un bootnode y 3 nodos adicionales
      1. Iniciar el Bootnode y obtener la dirección ENR necesaria para iniciar los nodos restantes
         1. Agregue el servicio a la red provista por el testsuite
            1. Implemente el objeto `services.ContainerCreationConfig` puede guiarse revisando como está configurado para el test `my_test` y sumando los cambios necesarios
               1. Utilice los archivos estáticos `genesis.json`, `password.txt` y`UTC--2021-08-11T21-30-29.861585000Z--14f6136b48b74b147926c9f24323d16c1e54a026`, utilizando los identificadores configurados en el método `Configure` del test, con la función `WithStaticFiles` del builder
               1. Establezca los puertos 8545 para el protocolo `tcp` y el puerto 30303 para el protocolo `tcp` y `udp` utilizando la función `WithUsedPorts()` del builder
            1. Implemente la función anónima que devuelve el objeto `services.ContainerRunConfig`
               1. Utilice la función `services.NewContainerRunConfigBuilder()` para crear el objeto `services.ContainerRunConfig` a retornar
                  1. Obtenga el filepath del archivo `genesis.json` utilizando el mapa `staticFileFilepaths` que recibe la función anónima como argumento
                  1. Obtenga el filepath del archivo `password.txt` utilizando el mapa `staticFileFilepaths` que recibe la función anónima como argumento
                  1. Obtenga el filepath del archivo `UTC--2021-08-11T21-30-29.861585000Z--14f6136b48b74b147926c9f24323d16c1e54a026` utilizando el mapa `staticFileFilepaths` que recibe la función anónima como argumento
                  1. Cree el entryPoint a ejecutar en el container con el comando necesario para lanzar un nodo de Ethereum,
                     1. Escriba el comando para inicializar el genesis block
                        1. Especifique un valor para el `datadir`
                        1. Especifique la ubicación del archivo de genesis block usando el filepath previamente obtenido
                     1. Escriba el comando para inicializar el nodo dentro de la red. El mismo deberá escribirlo a continuación del comando de inicialización y usando el operador '&&' para ejecutar comandos secunciales en un entrypoint
                        1. Especifique la ubicación del `datadir`
                        1. Especifique la ubicación del archivo `keystore`
                        1. Especifique el `network ID` recuerde que este valor ya está definido en el genesis block
                        1. Habilite el servidor `HTTP-RPC`
                        1. Especifique la `dirección IP` para el servidor HTTP-RPC
                        1. Especifique las `API's ofrecidas` sobre la interfaz HTTP-RPC colocandole los siguientes valores `admin,eth,net,web3,miner,personal,txpool,debug`
                        1. Especifique la `aceptación de todos los dominios` como origen colocando el valor '*'
                        1. Especifique la `IP` del nodo
                        1. Especifique el `puerto` que la red privada estará escuchando
                        1. Desbloquee la `cuenta firmante` para que pueda minar, recuerde que la dirección de esta cuenta la puede encontrar dentro del archivo keystore `UTC--2021-08-11T21-30-29.861585000Z--14f6136b48b74b147926c9f24323d16c1e54a026`
                        1. Habilite la `minería`
                        1. Habilite el `desbloqueo de cuenta inseguro` para poder desbloquear la cuenta del firmante
                        1. Restrinja la comunicación en la red configurando la `IP de la red` (CIDR masks)
                        1. Especifique el `filepath del archivo de password` que permite evitar la introducción del password manualmente cuando se quiere ejecutar un comando
            1. Llame a la función `networkCtx.AddService` pasandole un identificador del servicio y los 2 argumentos establecidos en los pasos anteriores
         1. Controle que el servicio esta correctamente cargado y funcionando
            1. Agregue las siguientes funciones privadas auxiliares `waitForStartup`, `isAvailable`, `getEnodeAddress` and `sendRpcCall` al final del archivo para que queden disponibles para utilizarlos luego
               ```
               func waitForStartup(ipAddress string, timeBetweenPolls time.Duration, maxNumRetries int) error {
                  for i := 0; i < maxNumRetries; i++ {
                     if isAvailable(ipAddress) {
                        return nil
                     }

                  // Don't wait if we're on the last iteration of the loop, since we'd be waiting unnecessarily
                  if i < maxNumRetries-1 {
                      time.Sleep(timeBetweenPolls)
                    }
                  }
                  return stacktrace.NewError(
                     "Service with ip '%v' did not become available despite polling %v times with %v between polls",
                     ipAddress,
                     maxNumRetries,
                     timeBetweenPolls)
               }

               func isAvailable(ipAddress string) bool {
                  enodeAddress, err := getEnodeAddress(ipAddress)
                  if err != nil {
                     return false
                  } else {
                     return strings.HasPrefix(enodeAddress, enodePrefix)
                  }
               }

               func getEnodeAddress(ipAddress string) (string, error) {
                  nodeInfoResponse := new(NodeInfoResponse)
                  err := sendRpcCall(ipAddress, adminInfoRpcCall, nodeInfoResponse)
                  if err != nil {
                     return "", stacktrace.Propagate(err, "Failed to send admin node info RPC request to geth node with ip %v", ipAddress)
                  }
                  return nodeInfoResponse.Result.Enode, nil
               }
               
               func sendRpcCall(ipAddress string, rpcJsonString string, targetStruct interface{}) error {
                  rpcPort := 8545
                  rpcRequestTimeout := 30 * time.Second
                  url := fmt.Sprintf("http://%v:%v", ipAddress, rpcPort)
                  var jsonByteArray = []byte(rpcJsonString)
               
                  logrus.Debugf("Sending RPC call to '%v' with JSON body '%v'...", url, rpcJsonString)
               
                  client := http.Client{
                     Timeout: rpcRequestTimeout,
                  }
                  resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonByteArray))
                  if err != nil {
                     return stacktrace.Propagate(err, "Failed to send RPC request to geth node with ip '%v'", ipAddress)
                  }
                  defer resp.Body.Close()
               
                  if resp.StatusCode == http.StatusOK {
                     // For debugging
                     var teeBuf bytes.Buffer
                     tee := io.TeeReader(resp.Body, &teeBuf)
                     bodyBytes, err := ioutil.ReadAll(tee)
                     if err != nil {
                        return stacktrace.Propagate(err, "Error parsing geth node response into bytes.")
                     }
                     bodyString := string(bodyBytes)
                     logrus.Tracef("Response for RPC call %v: %v", rpcJsonString, bodyString)
               
                     err = json.NewDecoder(&teeBuf).Decode(targetStruct)
                     if err != nil {
                        return stacktrace.Propagate(err, "Error parsing geth node response into target struct.")
                     }
                     return nil
                  } else {
                     return stacktrace.NewError("Received non-200 status code rom admin RPC api: %v", resp.StatusCode)
                  }
               }
               ```
            1. Implemente la llamada al método `waitForStartup` para controlar si el servicio previamente agregado está corriendo exitosamente
         1. Obtenga la dirección ENR del Bootnode 
            1. Ejecute el comando `geth attach data/geth.ipc --exec admin.nodeInfo.enr` pegue las siguientes líneas de código para realizar este paso
            ```
            exitCode, logOutput, err := serviceCtx.ExecCommand([]string{
               "/bin/sh",
               "-c",
               fmt.Sprintf("geth attach data/geth.ipc --exec admin.nodeInfo.enr"),
            })
            if err != nil {
               return "", stacktrace.Propagate(err, "Executing command returned an error with logs: %+v", string(*logOutput))
            }
            if exitCode != execCommandSuccessExitCode {
               return "", stacktrace.NewError("Executing command returned an failing exit code with logs: %+v", string(*logOutput))
            }
            ```
            2. Convierta el valor de logout del comando anterior en un string para obtener la dirección ENR del bootnode
      1. Iniciar el resto de los nodos con la ayuda del bootnode para poder conectar con la red
         1. Agregue el primer nodo
            1. Implemente el objeto `services.ContainerCreationConfig` 
               1. Utilice los archivos estáticos `genesis.json`y `password.txt` utilizando los identificadores configurados en el método `Configure` del test, con la función `WithStaticFiles` del builder
               1. Establezca los puertos `8545` para el protocolo `tcp` y el puerto `30303` para el protocolo `tcp` y `udp` utilizando la función `WithUsedPorts()` del builder
            1. Implemente la función anónima que devuelve el objeto `services.ContainerRunConfig`
               1. Utilice la función `services.NewContainerRunConfigBuilder()` para crear el objeto `services.ContainerRunConfig` a retornar
                  1. Obtenga el filepath del archivo `genesis.json` utilizando el mapa `staticFileFilepaths` que recibe la función anónima como argumento
                  1. Obtenga el filepath del archivo `password.txt` utilizando el mapa `staticFileFilepaths` que recibe la función anónima como argumento
                  1. Cree el entryPoint a ejecutar en el container con el comando necesario para lanzar un nodo de Ethereum,
                     1. Escriba el comando para inicializar el genesis block 
                        1. Especifique un valor para el `datadir`
                        1. Especifique la ubicación del archivo de genesis block usando el filepath previamente obtenido
                     1. Escriba el comando para inicializar el nodo dentro de la red. El mismo deberá escribirlo a continuación del comando de inicialización y usando el operador `&&` para ejecutar comandos secunciales en un entrypoint
                        1. Especifique la ubicación del `datadir`
                        1. Especifique el `network ID` recuerde que este valor ya está definido en el genesis block
                        1. Habilite el servidor `HTTP-RPC`
                        1. Especifique la `dirección IP` para el servidor HTTP-RPC
                        1. Especifique las `API's ofrecidas` sobre la interfaz HTTP-RPC colocandole los siguientes valores `admin,eth,net,web3,miner,personal,txpool,debug`
                        1. Especifique la `aceptación de todos los dominios` como origen colocando el valor '*'
                        1. Especifique la `IP` del nodo
                        1. Especifique el `puerto` que la red privada estará escuchando
                        1. Establezca el `bootnode` colocando la dirección ENR del bootnode
            1. Llame a la función `networkCtx.AddService` pasandole un identificador del servicio y los 2 argumentos establecidos en el paso anterior
         1. Controle que el nodo está correctamente cargado y funcionando
            1. Implemente nuevamente la llamada al método `waitForStartup` para controlar si el nodo está corriendo exitosamente
         1. Obtenga la dirección `enode` del nodo ya que el resto de los nodos la necesitarán para poder vincularse a este   
         1. Conecte el nodo con sus pares
            1. Conecte el nodo cargado uno por uno con el resto de los nodos (a excepción del bootnode) que están ejecutandose. Puede conectarlos manualmente utilizando el comando `admin_addPeer` [que esta explicado en esta documentación](https://geth.ethereum.org/docs/rpc/ns-admin#admin_addpeer)
            1. Compruebe que el nodo esta correctamente vinculado al resto. Puede listar la cantidad de peers utilizando el comando `admin_peers` [explicado en esta documentación](https://geth.ethereum.org/docs/rpc/ns-admin#admin_peers)
         1. Repita los pasos anteriores para iniciar y vincular los nodos restantes a la red
1. Implementar el método Run() del test
   1. Ejecute una transacción para deployar el smart contract `hello_world` dentro de la red privada
      1. Cree una instancia de Geth client que le será necesario para deployar el smart contract
      1. Obtenga el private key de la cuenta firmante
         1. Obtenga el contenido del archivo key de la cuenta firmante mediante el archivo `UTC--2021-08-11T21-30-29.861585000Z--14f6136b48b74b147926c9f24323d16c1e54a026` previamente cargado en el testsuite
         1. Obtenga el contenido del password mediante el archivo `password.txt` previamente cargado en el testsuite
         1. Obtenga el valor del private key mediante la función `keystore.DecryptKey` y pasandole los valores obtenidos del archivo key y password
      1. Obtenga un transactor que va a ser necesario para poder ejecutar el deploy, recomendamos utilizar la función `NewKeyedTransactorWithChainID`
      1. Ejecute la transacción para deployar el smart contract utilizando la función `DeployHelloWorld` provista por el binding `hello_world.go` la cual recibirá como parámetros el transactor y el Geth client
      1. Controle que la transacción ha terminado exitosamente
         1. Agregue el método auxiliar privado `waitUntilTransactionMined` al final del archivo
         ```
         func waitUntilTransactionMined(validatorClient *ethclient.Client, transactionHash common.Hash) error {
            maxNumCheckTransactionMinedRetries := 10
            timeBetweenCheckTransactionMinedRetries := 1 * time.Second
            for i := 0; i < maxNumCheckTransactionMinedRetries; i++ {
               receipt, err := validatorClient.TransactionReceipt(context.Background(), transactionHash)
               if err == nil && receipt != nil && receipt.BlockNumber != nil {
                  return nil
               }
               if i < maxNumCheckTransactionMinedRetries-1 {
                  time.Sleep(timeBetweenCheckTransactionMinedRetries)
               }
            }
            return stacktrace.NewError(
               "Transaction with hash '%v' wasn't mined even after checking %v times with %v between checks",
               transactionHash.Hex(),
               maxNumCheckTransactionMinedRetries,
               timeBetweenCheckTransactionMinedRetries)
            }
         ```
         1. Utilice el método auxiliar privado `waitUntilTransactionMined` para controlar si la transacción se realizó exitosamente
      1. Controle el funcionamiento del smart contract `HelloWorld` para esto puede utilizar la función `helloWorld.Greet()`
   1. Valide que el bootnode contiene la cuenta firmante, puede utilizar el comando `eth.accounts` para obtener la lista de cuentas
1. Agregue el test, dentro del archivo `my_testsuite.go`, a la lista de test a ejecutar cuando se lance el testsuite
1. Verify that running `bash scripts/build-and-run.sh all` generates output indicating that two test ran (my_test and my_advanced_test) and both are successfully executed





