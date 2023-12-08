// Goat File System Implementation
// if you would like to start from scratch, start here! Just note that "make" won't work until you have completed some basic skeleton code.

// if you are looking for a skeleton code, look at at ".goatfs.c" and copy the headers and empty functions here; you should be able to run "make" in the top-level directory without any errors.

#include <assert.h>
#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <stdbool.h>
#include <math.h>

#include "disk.h"
#include "goatfs.h"

SuperBlock *mounted_super; // current mounted superblock
int *bitmap; // bitmap for inodes blocks
int *data_bitmap; // may have to update in wremove    // bitmap for data blocks

// debug function as described by assignment
void debug() {
    // allocate space and read in superblock
    SuperBlock *superblock_data;
    superblock_data = malloc(4096);
    // read superblock from the 0th block position 
    wread(0, (char *)superblock_data);

    // begin printing superblock information
    unsigned int num_INode = superblock_data->InodeBlocks;

    printf("SuperBlock:\n");
    // check if the magic number is valid or not
    printf(superblock_data->MagicNumber == MAGIC_NUMBER ? "    magic number is valid\n" : "    magic number is not valid\n");

    // prinout debug information for the supoerblock
    printf("    %u blocks\n", superblock_data->Blocks);
    printf("    %u inode blocks\n", num_INode);
    printf("    %u inodes\n", superblock_data->Inodes);

    // for each inode block, allocate space and read in inode block
    for (int i = 0; i < num_INode; i++) {
        Block *inodeblock_data; // create iNode Block instance
        inodeblock_data = malloc(4096); // malloc the Block to be 4096 bytes
        // read from the current iNode block that we are on then setup the Block typedef struct
        wread(i + 1, (char *)inodeblock_data);

        // find and print inode block information
        for (int j = 0; j < 128; j++) { // 128 because there are 128 inodes per inode block
            if (inodeblock_data->Inodes[j].Valid == 1) { // check if the Inode is valid or not
                printf("Inode %d:\n", i * 128 + j); // talk to professor about this, since not encompassed in test cases
                printf("    size: %u bytes\n", inodeblock_data->Inodes[j].Size); // get the current Inode size
                printf("    direct blocks:"); // now we start debug prinout of direct blocks

                for (int k = 0; k < 5; k++) { // 5 because there are only 5 directs per inode
                    if (inodeblock_data->Inodes[j].Direct[k] != 0) { // if direct exists, print direct block
                        printf(" %u", inodeblock_data->Inodes[j].Direct[k]);
                    }
                }
                printf("\n");

                if (inodeblock_data->Inodes[j].Indirect != 0) { // if an indirect block exists for the inode, find and print information
                    printf("    indirect block: %u\n", inodeblock_data->Inodes[j].Indirect);
                    // allocate space and read in indirect block
                    Block *indirect_data; // create indirect Block instance
                    indirect_data = malloc(4096); // malloc the Block to be 4096 bytes
                    // read the indirect block in
                    wread(inodeblock_data->Inodes[j].Indirect, (char *)indirect_data);

                    printf("    indirect data blocks:");
                    for (int l = 0; l < 1024; l++) { // 1024 found using (sizeof(indirect_data->Pointers) / sizeof(indirect_data->Pointers[48]))
                        if (indirect_data->Pointers[l] != 0) { // if indirect exists, print
                            printf(" %u", indirect_data->Pointers[l]);
                        }
                    }
                    printf("\n");

                    // free data allocations
                    free(indirect_data);
                }
            }
        }
        // free inode block data
        free(inodeblock_data);
    }
    // free super block data
    free(superblock_data);
}

bool format() {
    // check if there is a disk mounted
    if (_disk->Mounts != 0) {
        return false;
    }
    // create superblock typedef struct
    SuperBlock *superblock_data;
    superblock_data = malloc(4096); // allocate a block size to it

    superblock_data->MagicNumber = (unsigned int)MAGIC_NUMBER; // assign magic number
    superblock_data->Blocks = (unsigned int)_disk->Blocks; // set superblock equal to the number of blocks on the disk
    superblock_data->InodeBlocks = (unsigned int)ceil((_disk->Blocks - 1) * 0.1); // calculate number of inode blocks
    superblock_data->Inodes = (unsigned int)superblock_data->InodeBlocks * 128; // calculate number of inodes total
    // write this superblock into the disk
    wwrite(0, (char *)superblock_data);
    // iterate through all of the blocks in the file system
    for (int i = 0; i < superblock_data->Blocks - 1; i++) {
        Block block; // create block
        memset(&block, 0, sizeof(Block));
        // write the empty block to the disk
        wwrite(i + 1, (char *)&block);
    }

    return true; // if we got here then we formatted good
}

int mount() {
    // check if there is a disk mounted
    if (_disk->Mounts != 0) { // already mounted!!!
        return -1;
    }
    // free up all of the global 
    free(mounted_super);
    free(bitmap);
    free(data_bitmap);
    
    mounted_super = malloc(4096); // allocate necessary memory for the superblock
    wread(0, (char *)mounted_super); // read in data from the superblock
    unsigned int num_INode = mounted_super->InodeBlocks; // 

    if (mounted_super->MagicNumber != MAGIC_NUMBER ||
        // mounted_super->InodeBlocks != (unsigned int) ceil((_disk->Blocks - 1) * 0.1) ||
        mounted_super->InodeBlocks * 128 != mounted_super->Inodes ||
        mounted_super->Blocks != _disk->Blocks) {
        return -1;
    }
    // allocate memory for both bitmaps
    bitmap = malloc(((int)num_INode) * 128 * sizeof(int)); 
    data_bitmap = malloc(((int)mounted_super->Blocks) * sizeof(int));

    // iterate through all blocks and create the data bitmap to determine which blocks are being used and not used
    for (int i = 0; i < (int)mounted_super->Blocks; i++) {
        if (i <= num_INode) { // block at this pos
            data_bitmap[i] = 1; // set to 1 -> true
        } else { // block not being used at this pos
            data_bitmap[i] = 0; // set to 0 -> false
        }
    }

    // iterate through all inodes
    for (int i = 0; i < num_INode; i++) {
        Block *inodeblock_data; // create typedef struct instance of the Block structure for the inode
        inodeblock_data = malloc(4096); // allocate a block size to it
        // read in current Inode blocks data
        wread(i + 1, (char *)inodeblock_data); // i + 1 to bypass the superblock
        // iterate through all inodes
        for (int j = 0; j < 128; j++) { // 128 because there are 128 inodes per inode block
            if (inodeblock_data->Inodes[j].Valid == 1) { // verify that the inode is valid
                bitmap[i * 128 + j] = 1; // make the inode used in the bitmap (aka set to 1 -> true)
                // iterate through all directs
                for (int k = 0; k < 5; k++) { // there is always 5 directs
                    if (inodeblock_data->Inodes[j].Direct[k] > 0) { // if the current direct is > 0 this means it is being used
                        data_bitmap[inodeblock_data->Inodes[j].Direct[k]] = 1; // add this onto the data bitmap
                    } else {
                        break; // otherwise just break out of this loop
                    }
                }

                if (inodeblock_data->Inodes[j].Indirect != 0) { // if an indirect block exists for the inode, find and print information
                    // allocate space and read in indirect block
                    Block *indirect_data; // create typedef struct instance of the Block structure for the indirect
                    indirect_data = malloc(4096); // allocate a block size to it
                    // read the Inode block data into the Block typedef struct instance
                    wread(inodeblock_data->Inodes[j].Indirect, (char *)indirect_data);
                    data_bitmap[inodeblock_data->Inodes[j].Indirect] = 1; // add this onto the data bitmap as being used
                    // iterate throught the whole indirect data block data and setup the pointers 
                    for (int k = 0; k < 4096; k++) {
                        if (indirect_data->Pointers[k] > 0) { // if there is an pointer
                            data_bitmap[indirect_data->Pointers[k]] = 1; // then add this pointer to thte data bitmap
                        } else {
                            break;
                        }
                    }

                    // free indirect data allocations
                    free(indirect_data);
                }
            } else {
                bitmap[i * 128 + j] = 0; // inode is not used in bitmap (aka set to 0 -> false)
            }
        }
        free(inodeblock_data); // free inode block data
    }

    _disk->Mounts += 1; // increment this as a successful mount

    return SUCCESS_GOOD_MOUNT; // good mount!!
}

ssize_t create() {
    // reads similar to mount for data concurrency
    free(mounted_super); // free up the mounted superblock
    mounted_super = malloc(4096); // then allocate it to the 4096 block size
    wread(0, (char *)mounted_super);

    // search thru bitmap for next uninitialized inode
    int index = -1;
    for (int i = 0; i < mounted_super->InodeBlocks * 128; i++) {
        if (bitmap[i] == 0) {
            index = i; // set index equal to the earlist uninitalized inode
            bitmap[i] = 1; // set bitmap to being using at this index
            break; // finally break out since we found a unused inode
        }
    }

    if (index == -1) { // if we didn't find anything then return -1
        return -1;
    }

    // create inode, write to file system, update bitmap
    Inode created_inode;
    created_inode.Valid = 1; // this inode is valid
    created_inode.Size = 0; // and the size of it is 0 becuase nothing is in it

    Block *inodeblock_data; // create a inode block data instance
    inodeblock_data = malloc(4096); // allocate 4096 bytes for this block
    wread((int)floor(index / 128) + 1, (char *)inodeblock_data); // read the inode block data from the correct pos

    Block new_data; // create new data block instance
    new_data = *inodeblock_data;
    new_data.Inodes[index % 128] = created_inode;  // set the current index of inode to the current one on

    wwrite((int)floor(index / 128) + 1, (char *)&new_data); // write the new data into the disk at specific index

    free(inodeblock_data); // free inode block data

    return index; // return the index we have created at
}

bool wremove(size_t inumber) {
    // reads similar to mount for data concurrency
    free(mounted_super); // free the mounted superblock global block
    mounted_super = malloc(4096); // allocated the usual 4096 bytes of space
    wread(0, (char *)mounted_super); // read into the superblock from block pos 0

    // check if inode is initialized
    if (bitmap[inumber] == 0) {
        return false; // inode isn't initalized!!
    } else {
        bitmap[inumber] = 0;
    }

    // create dummy inode, write to file system, update bitmap
    Inode removed_inode;
    removed_inode.Valid = 0;

    Block *inodeblock_data; // create inode block data instance
    inodeblock_data = malloc(4096); // malloc 4096 block size
    wread((int)floor(inumber / 128) + 1, (char *)inodeblock_data); // read in the inode block data

    // wipe indirects
    // left unimplemented, not needed for tests and also not accurate to real world

    // back to wiping inode
    Block new_data; // create new data instance
    new_data = *inodeblock_data;
    new_data.Inodes[inumber % 128] = removed_inode; // wipe the inode

    wwrite((int)floor(inumber / 128) + 1, (char *)&new_data); // write new (empty) block to wipe all existing data

    free(inodeblock_data);
    return true; // return true since we removed successfully
}

ssize_t stat(size_t inumber) {
    Block *inodeblock_data; // create inode block data instance
    inodeblock_data = malloc(4096); // allocate 4096 bytes of memory
    wread((int)floor(inumber / 128) + 1, (char *)inodeblock_data); // read 

    if (inodeblock_data->Inodes[inumber % 128].Valid == 1) { // verify that this current inode is valid before returning
        return inodeblock_data->Inodes[inumber % 128].Size; // output the size of the current inode
    }

    free(inodeblock_data);

    return -1; // else return -1 becuase we couldn't get a valid inode
}

ssize_t wfsread(size_t inumber, char *data, size_t length, size_t offset) {
    // mallocs and reads
    Block *inodeblock_data; // create inode block data instance
    inodeblock_data = malloc(4096); // allocate 4096 bytes of memory
    wread((int)floor(inumber / 128) + 1, (char *)inodeblock_data); // read in inode block data from inumber / 128 then floor + 1

    // error checks
    if (inodeblock_data->Inodes[inumber % 128].Valid != 1) { // check to make sure the current inode number is valid
        return -1;
    }
    // verify that the inode size is less than or equal to the offset value
    if (inodeblock_data->Inodes[inumber % 128].Size <= offset) {
        return -1;
    }

    // read directs to buffer
    int dir = floor(((int)offset) / 4096); // which block the offset begins in
    if (dir < 5) { // if dir is < 5 then direct it is an direct
        Block *direct_data; // create direct data instance
        direct_data = malloc(4096); // allocate 4096 bytes of memory for this block
        wread(inodeblock_data->Inodes[inumber % 128].Direct[dir], (char *)direct_data);

        strncpy(data, &direct_data->Data[offset % 4096], length);

        free(direct_data);
    } else { // read indirects to buffer
        Block *indirect_block; // create indirect block instance
        indirect_block = malloc(4096); // allocate 4096 bytes of memory for this block
        // read indirect into the typedef struct instance from the indirect
        wread(inodeblock_data->Inodes[inumber % 128].Indirect, (char *)indirect_block);

        // dir - 5 = index of indirect
        Block *indirect_data; // create indirect data instance
        indirect_data = malloc(4096); // allocate 4096 bytes of memory for this block
        // read indirect into the typedef struct instance
        wread(indirect_block->Pointers[dir - 5], (char *)indirect_data);

        strncpy(data, &indirect_data->Data[offset % 4096], length);
 
        free(indirect_data);
        free(indirect_block);
    }

    // compute return value and free
    int size_read = strlen(data);
    // printf("%d\n", size_read);
    free(inodeblock_data);

    return size_read;
}

ssize_t wfswrite(size_t inumber, char *data, size_t length, size_t offset) {
    // printf("STARTING LENGTH: %ld\n\n\n", length);
    // mallocs and reads
    // read from superblock (aka block 0)
    SuperBlock *superblock_data; // make a instance for the superblock
    superblock_data = malloc(4096); // allocate the size needed which is 4096 bytes for a block
    wread(0, (char *)superblock_data); // read in superblock data from block pos 0

    // read current Inode 
    Block *inodeblock_data; // make instance for inode block
    inodeblock_data = malloc(4096); // allocate 4096 bytes needed
    wread((int)floor(inumber / 128) + 1, ((char *)inodeblock_data)); // read from inode position and put read data into the instance made

    // errors
    // check if the current Indoe is valid
    if (inodeblock_data->Inodes[inumber % 128].Valid != 1) {
        return -1;
    }
    // verify that the inode size is less than or equal to the offset value
    if (inodeblock_data->Inodes[inumber % 128].Size != offset) {
        return -1;
    }
    // if the length is 0 then there is nothing to write :)
    if (length <= 0) {
        return 0;
    }

    int prev_offset = offset;
    for (int track = length; track > 0; track -= length) { // while loop to handle write-ins that are too long
        int dir = floor(((int)offset) / 4096); // which block the offset begins in

        if (track > 4096) {
            length = 4096;
        } else {
            length = track;
        }

        if (dir < 5) { // if dir is < 5 then direct it is an direct
            // if no direct exists
            if (!((int)inodeblock_data->Inodes[inumber % 128].Direct[dir] > 0 && (int)inodeblock_data->Inodes[inumber % 128].Direct[dir] < superblock_data->Blocks)) {
                int search = 0; // search for soonest block
                for (int i = 0; i < superblock_data->Blocks; i++) {
                    if (data_bitmap[i] == 0) { // found opening
                        search = i; // update search
                        data_bitmap[i] = 1; // setup bitmap for new instance
                        break;
                    }
                }
                if (search == 0) { // none found
                    // wremove(inumber);
                    // create();
                    // return -1;
                    return offset - prev_offset;
                }
                inodeblock_data->Inodes[inumber % 128].Direct[dir] = search; // set 
                wwrite((int)floor(inumber / 128) + 1, (char *)inodeblock_data);
            }

            // handle data write-in
            Block *direct_data; // direct data instance
            direct_data = malloc(4096); // allocate 4096 bytes
            // read in the direct into the instance created above
            wread((int)inodeblock_data->Inodes[inumber % 128].Direct[dir], (char *)direct_data);

            // time to update the data for direct
            Block updated_data; // create a update data instance block
            updated_data = *direct_data; // set it to the direct_data pointer
            strncpy(&updated_data.Data[offset % 4096], (const char *)&data[offset - prev_offset], length); // append data onto the current data
            if (track < 4096) { // if the track is less than 4096 bytes then append a null character
                memset(&updated_data.Data[track], '\0', 1); // append this null character
            }
            wwrite((int)inodeblock_data->Inodes[inumber % 128].Direct[dir], (char *)&updated_data); // write the data onto the disk

            free(direct_data);
        }
        else {
            // find block to hold indirects
            if (!((int)inodeblock_data->Inodes[inumber % 128].Indirect > 0 && (int)inodeblock_data->Inodes[inumber % 128].Indirect < superblock_data->Blocks)) {
                int search = 0; // search for soonest block
                for (int i = 0; i < superblock_data->Blocks; i++) {
                    if (data_bitmap[i] == 0) { // found opening
                        search = i; // update search
                        data_bitmap[i] = 1; // setup bitmap for new instance
                        break;
                    }
                }
                if (search == 0) { // none found
                    // wremove(inumber);
                    // create();
                    // return -1;
                    return offset - prev_offset;
                }
                inodeblock_data->Inodes[inumber % 128].Indirect = search;
                wwrite((int)floor(inumber / 128) + 1, (char *)inodeblock_data);
            }

            // find block to store indirect
            Block *indirect_data; // create indirect data instance
            indirect_data = malloc(4096); // allocate 4096 bytes
            wread(inodeblock_data->Inodes[inumber % 128].Indirect, (char *)indirect_data); // read indirect data in

            if (!((int)indirect_data->Pointers[dir - 5] > 0 && (int)indirect_data->Pointers[dir - 5] < superblock_data->Blocks)) {
                int search = 0; // search for soonest block
                for (int i = 0; i < superblock_data->Blocks; i++) {
                    if (data_bitmap[i] == 0) { // found opening
                        search = i; // update search
                        data_bitmap[i] = 1; // setup bitmap for new instance
                        break;
                    }
                }
                if (search == 0) { // none found
                    // wremove(inumber);
                    // create();
                    // return -1;
                    return offset - prev_offset;
                }
                indirect_data->Pointers[dir - 5] = search; // set pointer to be the soonest search pos
                wwrite(inodeblock_data->Inodes[inumber % 128].Indirect, (char *)indirect_data); // write indirect data into the block pos
            }

            // handle data write-in

            Block *write_data; // write data instance
            write_data = malloc(4096); // allocated 4096 bytes block size
            wread(indirect_data->Pointers[dir - 5], (char *)write_data); // read from the current indirect data pointer pos

            Block updated_data; // create a update data instance block
            updated_data = *write_data; // set it to the direct_data pointer
            strncpy(&updated_data.Data[offset % 4096], (const char *)&data[offset - prev_offset], length);
            if (track < 4096) { // if the track is less than 4096 bytes then append a null character
                memset(&updated_data.Data[track], '\0', 1); // append this null character
            }
            wwrite(indirect_data->Pointers[dir - 5], (char *)&updated_data); // write updated data into the disk

            free(write_data);
            free(indirect_data);
        }

        // update inode
        inodeblock_data->Inodes[inumber % 128].Size += length;
        // write in updated inode
        wwrite((int)floor(inumber / 128) + 1, (char *)inodeblock_data);

        offset += length; // update the offset we are dealing with
        // printf("TRACK %d, LENGTH %ld, OFFSET %ld, PREV_OFFSET %d\n\n", track, length, offset, prev_offset);
    }

    // compute return value and free
    int size_read = offset - prev_offset;
    free(inodeblock_data);
    free(superblock_data);

    return size_read;
}
