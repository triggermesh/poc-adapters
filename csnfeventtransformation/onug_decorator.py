import logging
import requests
import json

class onug: 
    
    __finding = {}
    __finding["provider"] = {}
    __finding["source"] = {}
    __finding["event"] = {}
    __finding["resource"] = {}

    def __init__(self,url,raw_finding):
        self.__raw_finding = raw_finding
        self.__url = url
        self.__set_provider_from_raw_finding()
        self.__get_onug_mapping_for_provider()
        self.__set_finding_provider_data()
        self.__set_data_source_mapping_from_provider()
        self.__map_raw_finding_to_onug()
        return

    def __set_provider_from_raw_finding(self):
        # function needs support for multiple providers
        if "ocid1.compartment" in str(self.__raw_finding):
            self.__provider = "Oracle Cloud Infrastructure"
        
        if "Microsoft.Security" in str(self.__raw_finding):
            self.__provider = "Azure"
        
        logging.debug("__set_provider_from_raw_finding: Provider is: " + self.__provider)

    def __set_finding_provider_data(self):
        # function needs support for multiple providers
        self.__finding["provider"]["providerId"] = self.__provider_mapping["providerId"]
        logging.debug("__set_provider_from_finding: ProviderId is set to: " + self.__finding["provider"]["providerId"]) 
        self.__finding["provider"]["providerType"] = self.__provider_mapping["providerType"]
        logging.debug("__set_provider_from_finding: ProviderId is set to: " + self.__finding["provider"]["providerType"])
        self.__finding["provider"]["name"] = self.__provider_mapping["provider"]
        logging.debug("__set_provider_from_finding: Provider is set to: " + self.__finding["provider"]["name"])
        return

    def __get_onug_mapping_for_provider(self):
        logging.debug("__get_onug_mapping_for_provider: getting provider mapping.")
        try:
             self.__provider_mapping = requests.get(self.__url).json()[self.__provider]
        except requests.exceptions.HTTPError as err:
            logging.error("__get_onug_mapping_for_provider: FAILED TO GET REQUEST" + str(err))
            raise SystemExit(err)

    def __set_provider_mapping_from_provider(self):
        if self.__provider == "Oracle Cloud Infrastructure":
            oci_provider_file = open('../provider/oci_output.json')
            logging.debug("__set_provider_mapping_from_provider")
            self.__provider_mapping = json.loads(oci_provider_file.read())[self.__provider]
            oci_provider_file.close()
        return

    def __set_data_source_mapping_from_provider(self):
        # function needs support of multiple providers    
        if self.__provider == "Oracle Cloud Infrastructure":
            if "cloudguard" in self.__raw_finding["data"]["resourceId"]:
                self.__source = "Cloud Guard"
                logging.debug(" __set_data_source_mapping_from_provider: source is: " + self.__source)
                logging.debug("__set_data_source_mapping_from_provider: Source Name is: " + self.__provider_mapping["source"][self.__source]["sourceName"])
                self.__finding["source"]["sourceName"] = self.__provider_mapping["source"][self.__source]["sourceName"]
                logging.debug("__set_data_source_mapping_from_provider: Source id is: " + self.__provider_mapping["source"][self.__source]["sourceId"])
                self.__finding["source"]["sourceId"] = self.__provider_mapping["source"][self.__source]["sourceId"]

                self.__source_mapping = self.__provider_mapping["source"][self.__source]["alerts"]["__default__"]["alertMapping"]
                logging.debug("__set_data_source_mapping_from_provider: Default Alerting Set: " + str(self.__source_mapping))


                for alert_name in self.__provider_mapping["source"][self.__source]["alerts"]:
                    logging.debug("__set_data_source_mapping_from_provider: Alert Name is: " + alert_name)
                    if alert_name ==  self.__raw_finding["data"]["additionalDetails"]["problemName"]:
                        logging.debug("__set_data_source_mapping_from_provider: Found problem: " + alert_name)
                        for mapping in self.__provider_mapping["source"][self.__source]["alerts"][alert_name]["alertMapping"]:
                            logging.debug("__set_data_source_mapping_from_provider: Adding alert mapping: " + mapping)
                            self.__source_mapping[mapping] = self.__provider_mapping["source"][self.__source]["alerts"][alert_name]["alertMapping"][mapping]


    def __mapped_item_from_path(self, path, finding):
        logging.debug("get_mapped_item_from_path: path is: " + str(path))
        tuple_path = tuple(path.split("."))
        for key in tuple_path:
            finding = finding[key]
        logging.debug("get_mapped_item_from_path: mapped_item is: " + str(finding))
        return finding

    def __set_item_from_path(self, path, item):
        # Cheating
        logging.debug("__set_item_from_path: path is: " + str(path) + " Item is: " + str(item) )
        keys = path.split(".")
        self.__finding[keys[0]][keys[1]] = item

        

    def __map_raw_finding_to_onug(self):
        logging.debug("__map_raw_finding_to_onug: source mapping is: " )
        for element in self.__source_mapping:
            logging.debug("__map_raw_finding_to_onug: Mapped Value: " + str(type(self.__source_mapping[element]["mappedValue"])))
            if not(self.__source_mapping[element]["mappedValue"]):
                path = self.__source_mapping[element]["path"]
                logging.debug("map_to_onug: path is: " + str(path))
                #logging.debug("map_to_onug: finding is: " + str(finding))
                mapped_item = self.__mapped_item_from_path(path, self.__raw_finding)

                logging.debug(f'map_to_onug: {element} is mapped to {mapped_item} via path: {self.__source_mapping[element]["path"]}')
                self.__set_item_from_path(element, mapped_item)
            else:
                mapped_item = self.__source_mapping[element]["value"]
                logging.debug(f'map_to_onug: Static map_to_onug: {element} is mapped to {mapped_item} static: {self.__source_mapping[element]["value"]}')
                self.__set_item_from_path(element, mapped_item)



    def get_provider_data(self):
        return self.__finding["provider"]
    
    def get_finding(self):
        return self.__finding