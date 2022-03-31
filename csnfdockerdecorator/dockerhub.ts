/**
 * (C) Copyright IBM Corp. 2021.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import {CsnfEvent} from './types';
import axios from "axios";

const DOCKERHUB_API_ROOT = 'https://hub.docker.com/v2/repositories';

export default class DockerhubDecorator {

    async decorate(csnfEvent: CsnfEvent) {
        if (!csnfEvent ||
            !csnfEvent.resource) {
            return null;
        }

        const containerImage = csnfEvent.image;
        const containerImageComponents = containerImage.split('/');

        if (containerImageComponents.length === 1) {
            containerImageComponents.unshift('library');
            containerImageComponents.unshift('docker.io');
        }
        const [registry, namespace, imageWithTagName] = containerImageComponents;

        if (registry !== 'docker.io' || namespace !== 'library') {
            return null;
        }

        const [imageName, imageTag] = imageWithTagName.split(':');

        const [imageInfo, tagInfo] = await this.getImageInfo(namespace, imageName, imageTag)

        const decoration = {
            registry: registry,
            namespace: imageInfo.namespace,
            image: imageInfo.name,
            tag: imageTag,
            digests: tagInfo.images.map((image) => image.digest),
            imageLastUpdated: imageInfo.last_updated,
            tagLastUpdated: tagInfo.last_updated,
            description: imageInfo.description,
            starCount: imageInfo.star_count,
            pullCount: imageInfo.pull_count,
        };
        return decoration;
    }

    async getImageInfo(imageNamespace: string, imageName: string, imageTag: string) {

        const imageInfoReqUrl = `${DOCKERHUB_API_ROOT}/${imageNamespace}/${imageName}/`;
        // console.log(`getting image ${imageName} info from ${imageInfoReqUrl}`);
        const imageInfo = (await axios.get(imageInfoReqUrl)).data;

        const tagInfoReqUrl = `${imageInfoReqUrl}tags/${imageTag}/`;
        // console.log(`getting tag ${imageTag} info from ${tagInfoReqUrl}`);
        const tagInfo = (await axios.get(tagInfoReqUrl)).data;


        return [imageInfo, tagInfo];
    }
}
